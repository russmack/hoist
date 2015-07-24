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

            var topBtn = $('<button />')
                .attr({ class: 'ui small compact basic icon button' })
                .attr({ title: 'Top (list processes)' })
            ;
            $(td)
                .append(
                    $('<span />').attr({ style:'margin-left:5px' })
                .append(
                    $(topBtn)
                .append(
                    $('<i />').attr({ class:'browser icon' })
                )))
            ;

            $(topBtn).click( function() {
                $.getJSON('/containers/top/' + val, function() {
                        //console.log('requested');
                })
                .done(function(data) {
                    var tHead = document.getElementById('table-top').tHead.children[0];
                    var tBody = document.getElementById('table-top-body');
                    $('#tab-top-message').text('');
                    $('#tab-top').modal('show');
                    if (data.StatusCode === 500 || data.Titles === 'undefined') {
                        $(tHead).empty();
                        $(tBody).empty();
                        $('#tab-top-message').text('No processes are running.');
                        return;
                    }
                    $(tHead).empty();
                    //$('#tab-top #results').text(JSON.stringify(data));
                    for (var col=0; col<data.Titles.length; col++) {
                        var title = data.Titles[col];
                        var th = document.createElement('th');
                        var thVal = document.createTextNode(title);
                        th.appendChild(thVal);
                        tHead.appendChild(th);
                    }
                    $(tBody).empty();
                    for (var r=0; r<data.Processes.length; r++) {
                        var tr = tBody.insertRow();
                        var row = data.Processes[r];
                        for (var c=0; c<row.length; c++) {
                            var td = tr.insertCell();
                            var val = row[c];
                            td.innerHTML = val;
                        }
                    }
                })
                .fail(
                    function( jqxhr, textStatus, error ) {
                        var err = textStatus + ", " + error;
                        console.log( "Request Failed: " + err );
                    }
                )
            });

            var statsBtn = $('<button />').attr({ class:'ui small compact basic icon button' });
            /* Disabled until streaming implemented.
            $(td)
                .append(
                    $('<span />').attr({ style:'margin-left:5px' })
                .append(
                    $(statsBtn)
                .append(
                    $('<i />').attr({ class:'area chart icon' })
                )))
            ;
            */

            $(statsBtn).click( function() {
                $('#tab-stats').modal('show');
                $.getJSON('/containers/stats/' + val, function() {
                        //console.log('requested');
                })
                .done(function(data) {
                     $('#tab-stats #results').text(JSON.stringify(data));
                })
                .fail(
                    function( jqxhr, textStatus, error ) {
                        var err = textStatus + ", " + error;
                        console.log( "Request Failed: " + err );
                    }
                )
            });
            
            var changesBtn = $('<button />')
                .attr({ class: 'ui small compact basic icon button' })
                .attr({ title: 'Filesystem changes' })
                ;

            $(td)
                .append(
                    $('<span />').attr({ style:'margin-left:5px' })
                .append(
                    $(changesBtn)
                .append(
                    $('<i />').attr({ class:'write square icon' })
                )))
            ;

            $(changesBtn).click( function() {
                $('#tab-changes-message').text('');
                var t = document.getElementById('table-changes');
                var tbl = new Table();
                tbl.create(t);
                tbl.clear();

                $.getJSON('/containers/changes/' + val, function() {
                    //console.log('requested');
                })
                .done(function(data) {
                    if (data.length === 0) {
                        $('#tab-changes-message').text('No changes.');
                        $('#tab-changes').modal('show');
                        return;
                    }
                    var t = document.getElementById('table-changes');
                    var tbl = new Table();
                    tbl.create(t);
                    tbl.setHeader(['Kind', 'Path']);
                    tbl.addBody(data);
                    $('#tab-changes').modal('show');
                })
                .fail(
                    function( jqxhr, textStatus, error ) {
                        var err = textStatus + ", " + error;
                        console.log( "Request Failed: " + err );
                    }
                    )
            });

            
            var deleteBtn = $('<button />')
                .attr({ class: 'ui small compact basic icon button' })
                .attr({ title: 'Delete' })
            ;

            $(td)
                .append(
                    $('<span />').attr({ style:'margin-left:5px' })
                .append(
                    $(deleteBtn)
                .append(
                    $('<i />').attr({ class:'remove square icon' })
                )))
            ;

            $(deleteBtn).click( function() {
                $('#tab-delete').modal('show');
                $.getJSON('/containers/delete/' + val, function() {
                        //console.log('requested');
                })
                .done(function(data) {
                    var statusCode = data.StatusCode;
                    $('#tab-delete #results').text('Response status code: ' + statusCode);
                    $('#table-list-body').empty();
                    loadContainerList();
                })
                .fail(
                    function( jqxhr, textStatus, error ) {
                        var err = textStatus + ", " + error;
                        console.log( "Request Failed: " + err );
                    }
                )
            });

            /* Disable until post requests implemented.
            var startBtn = $('<button />').attr({ class:'ui small compact basic icon button' });
            $(td)
                .append(
                    $('<span />').attr({ style:'margin-left:5px' })
                .append(
                    $(startBtn)
                .append(
                    $('<i />').attr({ class:'play icon' })
                )))
            ;

            $(startBtn).click( function() {
                $('#tab-start').modal('show');
                $.getJSON('/containers/start/' + val, function() {
                        //console.log('requested');
                })
                .done(function(data) {
                     $('#tab-start #results').text(JSON.stringify(data));
                })
                .fail(
                    function( jqxhr, textStatus, error ) {
                        var err = textStatus + ", " + error;
                        console.log( "Request Failed: " + err );
                    }
                )
            });
            */

            /* Disabled until post requests implemented.
            var stopBtn = $('<button />').attr({ class:'ui small compact basic icon button' });
            $(td)
                .append(
                    $('<span />').attr({ style:'margin-left:5px' })
                .append(
                    $(stopBtn)
                .append(
                    $('<i />').attr({ class:'stop icon' })
                )))
            ;

            $(stopBtn).click( function() {
                $('#tab-stop').modal('show');
                $.getJSON('/containers/stop/' + val, function() {
                        //console.log('requested');
                })
                .done(function(data) {
                     $('#tab-stop #results').text(JSON.stringify(data));
                })
                .fail(
                    function( jqxhr, textStatus, error ) {
                        var err = textStatus + ", " + error;
                        console.log( "Request Failed: " + err );
                    }
                )
            });
            */

            var historyBtn = $('<button />').attr({ class:'ui small compact basic icon button' });
            $(historyBtn).click( function() {
                //$('#menu-tabs #menu-tab-history').click();
                $('#tab-history').modal('show');
                $.getJSON('/containers/log/' + val, function() {
                        //console.log('requested');
                })
                .done(function(data) {
                     $('#tab-history #results').text(JSON.stringify(data));
                })
                .fail(
                    function( jqxhr, textStatus, error ) {
                        var err = textStatus + ", " + error;
                        console.log( "Request Failed: " + err );
                    }
                )
            });


            /* Disabled containers logs button, until stream response understood.
            $(td)
                .append(
                    $('<span />').attr({ style:'margin-left:5px' })
                .append(
                    $(historyBtn)
                .append(
                    $('<i />').attr({ class:'history icon' })
                )))
            ;
            */

            $(link).click( function() {
                    $('#tab-inspect').modal('show');
                    $.getJSON('/containers/inspect/' + val, function() {
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

    loadContainerList();

    function loadContainerList() {
        $.getJSON("/containers/list")
            .done(function(data) {
                tableCreate($("#results")[0], data);
            })
        ;
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
