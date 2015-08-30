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

    var apiVersion = '0.1';

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
/*
            // Add container for buttons.
            var btns = $('<span />');

            // Add images button.
            var nodesBtn = buildButton('Nodes', 'cube icon');
            var link = document.createElement('a');
            //var imagesUri = 'nodes.html?clusterid=' + clusterId  + '&nodeid=' + idLink;
            var nodesUri = 'nodes.html?clusterid=' + idLink;
            console.log('URI: ' + nodesUri);
            link.setAttribute('href', nodesUri)
            $(link).append(nodesBtn);
            $(btns).append($(link));

            // Add containers button.
            var containersBtn = buildButton('Containers', 'cubes icon');
            var link = document.createElement('a');
            var containersUri = 'containers.html?clusterid=' + clusterId + '&nodeid=' + idLink;
            link.setAttribute('href', containersUri)
            $(link).append(containersBtn);
            $(btns).append($(link));
            $(td).append($(btns));
*/
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
            var td = $('td:eq(1)', this)[0];
            var val = td.childNodes[0].textContent;
            var btns = td.childNodes[1];
            var abbr = val.substring(0, 11);
            var link = document.createElement('a');
            var linkText = document.createTextNode(abbr);
            link.setAttribute('href', 'nodes.html?clusterid=' + $('td:eq(0)', this)[0].textContent);
            link.className = '';
            link.appendChild(linkText);
            td.innerHTML = '';  // Clear cell first.
            td.appendChild(link);
            /*
            td.appendChild(btns);
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

    loadClusterList();

    function loadClusterList() {
        $.getJSON('/' + apiVersion  + '/clusters')
            .done(function(data) {
                if (data.length === 0) {
                    $('#tab-list-message').text('No clusters.');
                    return;
                }
                var t = document.getElementById('table-list');
                var tbl = new Table();
                tbl.create(t);
                tbl.setHeader([
                    'Id', 
                    'Name', 
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

    var addClusterValidationRules = {
        'tab-add-name-text': {
            identifier: 'tab-add-name-text', 
            rules: [
            {
                type: 'empty',
                prompt: 'Give the cluster a name.'
            }
            ]
        },
        'tab-add-description-text': {
            identifier: 'tab-add-description-text', 
            rules: [
            {
                type: 'empty',
                prompt: 'Give the cluster a description.'
            }
            ]
        }

    };

    $('#tab-add .ui.form').form(addClusterValidationRules, { inline: true,  onSuccess: function() {
            $('.ui.dimmer').dimmer('show');
            
            var name = $('#tab-add-name-text').val();
            var desc = $('#tab-add-description-text').val();
            var cluster = {
                'Name': name,
                'Description': desc
            };
            var postBody = JSON.stringify(cluster);
            $.post('/' + apiVersion  + '/clusters', postBody, function() {
                console.log('success') } )
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
                    loadClusterList();
                    $('#tab-add-added').modal('show');
                })
                .fail(function(jqxhr, textStatus, error) {
                    var err = textStatus + ", " + error;
                    console.log('Request Failed: ' + err);
                })
            ;
            
        }})
    ;

    /*
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
    */

  })
;
