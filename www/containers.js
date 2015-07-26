$(document)
  .ready(function() {

    var
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
                    $('#tab-top-message').text('');
                    var tbl = new Table();
                    var t = document.getElementById('table-top');
                    tbl.create(t);
                    tbl.clear();
                    if (data.StatusCode === 500 || data.Titles === 'undefined') {
                        $('#tab-top-message').text('No processes are running.');
                        $('#tab-top').modal('show');
                        return;
                    }
                    tbl.setHeader(data.Titles);
                    tbl.addBody(data.Processes);
                    $('#tab-top').modal('show');
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
                //tableCreate($("#results")[0], data);
                    if (data.length === 0) {
                        $('#tab-list-message').text('No containers.');
                        //$('#tab-list').modal('show');
                        return;
                    }
                    var t = document.getElementById('table-list');
                    var tbl = new Table();
                    tbl.create(t);
                    tbl.setHeader([
                        'Id', 
                        'Image', 
                        'Command', 
                        'Created', 
                        'Status', 
                        'Ports', 
                        'Size Rw', 
                        'Size RootFs'
                    ]);
                    //tbl.addBody(data);
                    //$('#tab-list').modal('show');
                    tableCreate($("#results")[0], data);
            })
        ;
    }


  })
;
