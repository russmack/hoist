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
        $('#table-list > tbody > tr').each( function() {
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

            /*
             * Add buttons.
             */

            var historyBtn = buildButton('History', 'history icon');
            $(td).append($(historyBtn));

            var deleteBtn = buildButton('Delete', 'delete icon');
            $(td).append($(deleteBtn));


            /*
             * Bind button event handlers.
             */

            $(historyBtn).click( function() {
                $('#tab-history-message').text('');
                var t = document.getElementById('table-history');
                var tbl = new Table();
                tbl.create(t);
                tbl.clear();

                // Get inspect data.
                $.getJSON('/' + apiVersion + '/nodes/' + nodeId + '/images/history/' + val, function() {
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

            $(link).click( function() {
                    //$('#menu-tabs #menu-tabs-inspect').click();
                    // Get inspect data.
                    $.getJSON('/' + apiVersion + '/nodes/' + nodeId  + '/images/inspect/' + val, function() {
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

            $(deleteBtn).click( function() {
                $('#tab-delete').modal('show');
                $.getJSON('/' + apiVersion + '/nodes/' + nodeId  + '/images/delete/' + val, function() {
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
            
        });

    }

    $('#menu-tabs .item')
        .tab('change tab', 'tab-list')
    ;

    function loadImageList() {
        $.getJSON('/' + apiVersion + '/nodes/' + nodeId  + '/images/list')
            .done(function(data) {
                tableCreate($("#results")[0], data);
            })
        ;
    }

    var apiVersion = '0.1';

    var clusterId = document.getElementById('hidden-clusterid').value;
    var nodeId = document.getElementById('hidden-nodeid').value;

    loadImageList();

    $('#menu-tabs-search')
        .on('click', function() {
            $('#tab-search')
                .tab('change tab', 'tab-search')
            ;
        })
    ;

    var searchValidationRules = {
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

    $('.ui.form').form(searchValidationRules, { inline: true,  onSuccess: function() {
            //$('.ui.dimmer').dimmer('show');
            var term = $('#tab-search-text').val();
            $.getJSON('/' + apiVersion + '/images/search/' + term)
                .done(function(data) {
                    //$('.ui.dimmer').dimmer('hide');
                    $('#tab-search-results #results').text(renderSearchResults(data));
                    $('#tab-search-results').modal('show');
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
