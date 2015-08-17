$(document)
  .ready(function() {

    $('.ui.dropdown')
      .dropdown({
        on: 'hover'
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
            //var id = data[i].Id.substring(0, 11);
            //var id = data[i].Id;
            var idLink = data[i].Id;
            //var idLink = '<a href="images/inspect/' + data[i].Id  + '">' + id  + '</a>';
            //var idLink = '<a onclick=""' + data[i].Id  + '">' + id  + '</a>';
            //td.appendChild(document.createTextNode(idLink));
            td.innerHTML = idLink;
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Name));

            // Add images button.
            var imagesBtn = buildButton('Images', 'cube icon');
            var link = document.createElement('a');
            var imagesUri = 'images.html?nodeid=' + idLink;
            console.log('URI: ' + imagesUri);
            link.setAttribute('href', imagesUri)
            $(link).append(imagesBtn);
            $(td).append($(link));

            // Add containers button.
            var containersBtn = buildButton('Containers', 'cubes icon');
            var link = document.createElement('a');
            var containersUri = 'containers.html?nodeid=' + idLink;
            console.log('URI: ' + containersUri);
            link.setAttribute('href', containersUri)
            $(link).append(containersBtn);
            $(td).append($(link));

            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Scheme));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Address));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Port));
            var td = tr.insertCell();
            td.appendChild(document.createTextNode(data[i].Description));
            var td = tr.insertCell();
            //var dt = convertUnixTime(data[i].Created);
            var dt = data[i].Created;
            td.appendChild(document.createTextNode(dt));
            //var td = tr.insertCell();
            //var text = '';
            //if (data[i].Labels !== 'undefined' && data[i].Labels !== null)  {
            //    for (var j = 0; j < data[i].Labels.length; j++)
            //    {
            //        text += data[i].Labels[j] + '<br />';
            //    }
            //}
            //td.appendChild(document.createTextNode(text));
            //var td = tr.insertCell();
            //td.appendChild(document.createTextNode(data[i].RepoDigests));
        }
        $('#table-list > tbody > tr').each( function() {
            /* 
            var td = $('td:eq(0)', this)[0];
            var val = td.textContent;
            var abbr = val.substring(0, 11);
            var link = document.createElement('a');
            var linkText = document.createTextNode(abbr);
            link.setAttribute('href', '#')
            link.className = '';
            link.appendChild(linkText);
            td.innerHTML = '';  // Clear cell first.
            td.appendChild(link);
            */

            /*
             * Add buttons.
             */

            /*
            var historyBtn = buildButton('History', 'history icon');
            $(td).append($(historyBtn));

            var deleteBtn = buildButton('Delete', 'delete icon');
            $(td).append($(deleteBtn));
            */

            /*
             * Bind button event handlers.
             */

            /*
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
            */
            /*
            $(link).click( function() {
                    //$('#menu-tabs #menu-tabs-inspect').click();
                    // Get inspect data.
                    $.getJSON('/images/inspect/' + val, function() {
                            //console.log('requested');
                        })
                        .done(function(data) {
                            //$('#tab-inspect #results').text(JSON.stringify(data));
                            var html = renderJson(data);
                            $('#tab-inspect #results').html(html);
                            $('#tab-inspect').modal('show');
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
            */
            /*
            $(deleteBtn).click( function() {
                $('#tab-delete').modal('show');
                $.getJSON('/images/delete/' + val, function() {
                        //console.log('requested');
                })
                .done(function(data) {
                    var statusCode = data.StatusCode;
                    var html = renderJson(data);
                    $('#tab-delete #results').html(html);
                    $('#table-list-body').empty();
                    loadImageList();
                })
                .fail(
                    function( jqxhr, textStatus, error ) {
                        var err = textStatus + ", " + error;
                        console.log( "Request Failed: " + err );
                    }
                )
            });
            */
        });

    }

    $('#menu-tabs .item')
        .tab('change tab', 'tab-list')
    ;

    loadNodeList();

    function loadNodeList() {
        $.getJSON("/nodes/list")
            .done(function(data) {
                if (data.length === 0) {
                    $('#tab-list-message').text('No nodes.');
                    return;
                }
                var t = document.getElementById('table-list');
                var tbl = new Table();
                tbl.create(t);
                tbl.setHeader([
                    'Id', 
                    'Name', 
                    'Scheme',
                    'Address', 
                    'Port',
                    'Description', 
                    'Created', 
                    ]);
                tableCreate($("#tab-list #results")[0], data);
            })
        ;
    }


    /*
    $('#menu-tabs-search')
        .on('click', function() {
            $('#tab-search')
                .tab('change tab', 'tab-search')
            ;
        })
    ;

    $('#menu-tabs-add')
        .on('click', function() {
            $('#tab-add')
                .tab('change tab', 'tab-add')
            ;
        })
    ;
    */

    var addNodeValidationRules = {
        'tab-add-name-text': {
            identifier: 'tab-add-name-text', 
            rules: [
            {
                type: 'empty',
                prompt: 'Give the node a name.'
            }
            ]
        },
        'tab-add-scheme-text': {
            identifier: 'tab-add-scheme-text', 
            rules: [
            {
                type: 'empty',
                prompt: 'HTTP or HTTPS?'
            }
            ]
        },
        'tab-add-address-text': {
            identifier: 'tab-add-address-text', 
            rules: [
            {
                type: 'empty',
                prompt: 'What\'s the IP address?'
            }
            ]
        },
        'tab-add-port-text': {
            identifier: 'tab-add-port-text', 
            rules: [
            {
                type: 'empty',
                prompt: 'What\'s the port?'
            }, 
            {
                type: 'integer[1..65535]',
                prompt: 'Invalid port number.'
            }
            ]
        },
        'tab-add-description-text': {
            identifier: 'tab-add-description-text', 
            rules: [
            {
                type: 'empty',
                prompt: 'Give the node a description.'
            }
            ]
        }

    };

    $('#tab-add .ui.form').form(addNodeValidationRules, { inline: true,  onSuccess: function() {
            $('.ui.dimmer').dimmer('show');
            
            var name = $('#tab-add-name-text').val();
            var scheme = $('#tab-add-scheme-text').val();
            var address = $('#tab-add-address-text').val();
            var port = parseInt( $('#tab-add-port-text').val(), 10);
            var desc = $('#tab-add-description-text').val();
            var node = {
                'Name': name,
                'Scheme': scheme,
                'Address': address,
                'Port': port, 
                'Description': desc
            };
            var postBody = JSON.stringify(node);
            $.post('/nodes', postBody, function() { console.log('success') } )
                .done(function(data) {
                    $('.ui.dimmer').dimmer('hide');
                    var statusCode = data.StatusCode;
                    var html = renderJson(JSON.parse(data));
                    $('#tab-add-added #results').html(html);
                    $('#menu-tabs .item')
                        .tab('change tab', 'tab-list')
                    ;

                    var t = document.getElementById('table-list');
                    var tbl = new Table();
                    tbl.create(t);
                    tbl.clear();
                    loadNodeList();
                    $('#tab-add-added').modal('show');
                })
                .fail(function(jqxhr, textStatus, error) {
                    var err = textStatus + ", " + error;
                    console.log('Request Failed: ' + err);
                })
            ;
            
        }})
    ;

    function renderSearchResults(json) {
        $('#tab-search-results-message').text('');
        var t = document.getElementById('table-search-results');
        var tbl = new Table();
        tbl.create(t);
        tbl.clear();
        if (json.length === 0) {
            $('#tab-search-results-modal-message').modal('show');
            $('#tab-search-results-message').text('No results.');
            return;
        }
        tbl.setHeader(['Description', 'Is Official', 'Is Trusted', 'Name', 'Star Count']);
        tbl.addBody(json);
    }

  })
;
