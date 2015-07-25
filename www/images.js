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
        //var tbl  = document.createElement("table");
        //var tbl = document.getElementById('table-list');
        var listBody = $('#table-list > tbody')[0];

        for (var i = 0; i < data.length; i++)
        {
            var tr = listBody.insertRow();
            var td = tr.insertCell();
            var id = data[i].Id.substring(0, 11);
            var idLink = data[i].Id;
            //var idLink = '<a href="images/inspect/' + data[i].Id  + '">' + id  + '</a>';
            //var idLink = '<a onclick=""' + data[i].Id  + '">' + id  + '</a>';
            //td.appendChild(document.createTextNode(idLink));
            td.innerHTML = idLink;
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].RepoTags));
            var td = tr.insertCell();
            var dt = convertUnixTime(data[i].Created);
            td.appendChild(document.createTextNode(dt));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].VirtualSize));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].ParentId.substring(0, 11)));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Size));
            var td = tr.insertCell();
            var text = '';
            if (data[i].Labels !== 'undefined' && data[i].Labels !== null)  {
                for (var j = 0; j < data[i].Labels.length; j++)
                {
                    text += data[i].Labels[j] + '<br />';
                }
            }
            td.appendChild(document.createTextNode(text));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].RepoDigests));
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
            //var historyImage = 
            // $('<img />').attr({
            //       src:'some image url',
            //       width:'width in integer',
            //       height:'integer'
            // }).appendTo($('<a />').attr({
            //       href:'somelink'
            // }).appendTo($('#someElement')));

            var historyBtn = $('<button />')
                .attr({ class:'ui small compact basic icon button' })
                .attr({ title:'History' })
            ;
            $(historyBtn).click( function() {
                $('#tab-history-message').text('');
                var t = document.getElementById('table-history');
                var tbl = new Table();
                tbl.create(t);
                tbl.clear();

                // Get inspect data.
                $.getJSON('/images/history/' + val, function() {
                        //console.log('requested');
                })
                .done(function(data) {
                    if (data.length === 0) {
                        $('#tab-history-message').text('No history.');
                        $('#tab-history').modal('show');
                        return;
                    }
                    var t = document.getElementById('table-history');
                    var tbl = new Table();
                    tbl.create(t);
                    tbl.setHeader(['Created', 'Created By', 'Id', 'Size', 'Tags']);
                    tbl.addBody(data);
                    $('#tab-history').modal('show');
                })
            });

            $(td)
                .append(
                    $('<span />').attr({ style:'margin-left:5px' })
                .append(
                    $(historyBtn)
                .append(
                    $('<i />').attr({ class:'history icon' })
                )))
            ;

            $(link).click( function() {
                    //$('#menu-tabs #menu-tabs-inspect').click();
                    $('#tab-inspect').modal('show');
                    // Get inspect data.
                    $.getJSON('/images/inspect/' + val, function() {
                            //console.log('requested');
                        })
                        .done(function(data) {
                            $('#tab-inspect #results').text(JSON.stringify(data));
                        })
                        .fail(
                            function( jqxhr, textStatus, error ) {
                                var err = textStatus + ", " + error;
                                console.log( "Request Failed: " + err );
                            }
                        )
                        .always(function() {
                            //console.log( "complete" );
                        })
                    ;
                }
            );
        });

    }

    $('#menu-tabs .item')
        .tab('change tab', 'tab-list')
    ;

    $.getJSON("/images/list")
        .done(function(data) {
            tableCreate($("#results")[0], data);
        })
    ;

    $('#menu-tabs-search')
        .on('click', function() {
            $('#tab-search')
                .tab('change tab', 'tab-search')
            ;
        })
    ;

    var validationRules = {
        'tab-search-text': {
            identifier: 'tab-search-text', 
            rules: [
            {
                type: 'empty',
                prompt: 'What would you like to search for?'
            }
            ]
        }
    };

    $('.ui.form').form( validationRules , { inline: true,  onSuccess: function() {
            $('.ui.dimmer').dimmer('show');
            var term = $('#tab-search-text').val();
            $.getJSON('/images/search/' + term)
                .done(function(data) {
                    $('.ui.dimmer').dimmer('hide');
                    $('#tab-search-results').modal('show');
                    $('#tab-search-results #results').text(JSON.stringify(data));
                })
                .fail(function(jqxhr, textStatus, error) {
                    var err = textStatus + ", " + error;
                    console.log('Request Failed: ' + err);
                })
            ;
        }})
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

    function Table() {
        this.table;
    }
    Table.prototype.create = function(tableElem) {
        this.table = tableElem;
    };
    Table.prototype.setHeader = function(v) {
        var th = this.table.tHead;
        if (th !== null) {
            this.table.deleteTHead();
        }
        th = this.table.createTHead();
        var row = th.insertRow(0);
        for (var i=0; i<v.length; i++) {
            var cell = row.insertCell(i);
            var txt = document.createTextNode(v[i]);
            cell.appendChild(txt);
        }
    };
    Table.prototype.addBody = function(v) {
        var newtbody = document.createElement('tbody');
        var tbody = this.table.tBodies;
        if (tbody !== null) {
            tbody[0].parentNode.replaceChild(newtbody, tbody[0]);
        }
        var b = this.table.tBodies[0];
        for (var r=0; r<v.length; r++) {
            var row = b.insertRow();
            for (prop in v[r]) {
                if( v[r].hasOwnProperty(prop) ) {
                    var cell = row.insertCell();
                    var txt = document.createTextNode(v[r][prop]);
                    cell.appendChild(txt);
                }
            }
        }
    };
    Table.prototype.clearHeader = function() {
        if (this.table.tHead !== null) {
            this.table.deleteTHead();
        }
    };
    Table.prototype.clearBody = function() {
        //for (var i=0; i<this.table.rows
        var newtbody = document.createElement('tbody');
        var tbody = this.table.tBodies;
        if (tbody !== null) {
            tbody[0].parentNode.replaceChild(newtbody, tbody[0]);
        }
    };
    Table.prototype.clear = function() {
        this.clearHeader();
        this.clearBody();
    };

  })
;
