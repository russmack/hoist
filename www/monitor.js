$(document)
  .ready(function() {

    var
      changeSides = function() {
        $('.ui.shape')
          .eq(0)
            .shape('flip over')
            .end()
          .eq(1)
            .shape('flip over')
            .end()
          .eq(2)
            .shape('flip back')
            .end()
          .eq(3)
            .shape('flip back')
            .end()
        ;
      },
      validationRules = {
        firstName: {
          identifier  : 'email',
          rules: [
            {
              type   : 'empty',
              prompt : 'Please enter an e-mail'
            },
            {
              type   : 'email',
              prompt : 'Please enter a valid e-mail'
            }
          ]
        }
      }
    ;

    $('.ui.dropdown')
      .dropdown({
        on: 'hover'
      })
    ;

    $('.ui.form')
      .form(validationRules, {
        on: 'blur'
      })
    ;

    $('.masthead .information')
      .transition('scale in', 1000)
    ;

    // Function to create a table.
    function tableCreate(el, data)
    {
        var listBody = $('#table-list > tbody')[0];

        for (var i = 0; i < data.length; i++)
        {
            var tr = listBody.insertRow();
            var td = tr.insertCell();
            var id = data[i].Id.substring(0, 11);
            var idLink = data[i].Id;
            td.innerHTML = idLink;
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Image));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Command));
            var td = tr.insertCell();
            var dt = convertUnixTime(data[i].Created);
            td.appendChild(document.createTextNode(dt));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Status));
            var td = tr.insertCell();
            var text = '';
            if (data[i].Ports !== 'undefined' && data[i].Ports !== null)  {
                for (var j = 0; j < data[i].Ports.length; j++)
                {
                    text += JSON.stringify(data[i].Ports[j]);
                }
            }
            td.appendChild(document.createTextNode(text));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].SizeRw));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].SizeRootFs));
        }
        //el.appendChild(tbl);
        $('#table-list > tbody > tr').each( function() {
            var td = $('td:eq(0)', this)[0];
            var val = td.innerText;
            var abbr = val.substring(0, 11);
            var link = document.createElement('a');
            var linkText = document.createTextNode(abbr);
            link.setAttribute('href', '#')
            link.className = '';
            link.appendChild(linkText);
            td.innerHTML = '';  // Clear cell first.
            td.appendChild(link);
            var br = document.createElement('br');
            td.appendChild(br);





        });

    }

    $('#menu-tabs .item')
        .tab('change tab', 'tab-info')
    ;
    loadInfoTab();

    function loadInfoTab() {
        $.getJSON("/monitor/info")
            .done(function(data) {
                $('#tab-info #results').text(JSON.stringify(data));
            })
        .fail(function( jqxhr, textStatus, error ) {
            var err = textStatus + ", " + error;
            console.log( "Request Failed: " + err );
        })
        ;
    }

    $('#menu-tabs-info')
        .on('click', function() {
            loadInfoTab();
        })
    ;
    
    $('#menu-tabs-dockerversion')
        .on('click', function() {
            $.getJSON("/monitor/version")
                .done(function(data) {
                    $('#tab-dockerversion #results').text(JSON.stringify(data));
                })
                .fail(function( jqxhr, textStatus, error ) {
                    var err = textStatus + ", " + error;
                    console.log( "Request Failed: " + err );
                })
            ;
        })
    ;

    $('#menu-tabs-ping')
        .on('click', function() {
            $.getJSON("/monitor/ping")
                .done(function(data) {
                    $('#tab-ping #results').text(JSON.stringify(data));
                })
                .fail(function( jqxhr, textStatus, error ) {
                    var err = textStatus + ", " + error;
                    console.log( "Request Failed: " + err );
                })
            ;
        })
    ;

    var jsonStream;

    $('#menu-tabs-events')
        .on('click', function() {
            if (!!window.EventSource) {
               console.log('Event sourcing is available.');
            } else {
                console.log('Event sourcing is not available in this browser.');
            }
            
            jsonStream = new EventSource('monitor/events');
            
            jsonStream.addEventListener('message', function(e) {
                    console.log(e.data);
                    $('#tab-events #results').append(e.data + '<br />');
                }, false)
            ;
            
            jsonStream.addEventListener('open', function(e) {
                    console.log("opened channel")
                }, false)
            ;
           
            jsonStream.addEventListener('error', function(e) {
                    if (e.readyState == EventSource.CLOSED) {
                        console.log("closed channel")
                    } else {
                        console.log('error: ' + e);
                    }
                }, false)
            ;

        })
    ;

    function convertUnixTime(unix_timestamp) {
       var date = new Date(unix_timestamp*1000);
      // hours part from the timestamp
      var hours = date.getHours();
      // // minutes part from the timestamp
      var minutes = "0" + date.getMinutes();
      // // seconds part from the timestamp
      var seconds = "0" + date.getSeconds();
      
      // // will display time in 10:30:23 format
      var formattedTime = hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
      //return formattedTime;
      return date;
    }

    setInterval(changeSides, 3000);


  })
;
