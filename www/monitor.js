$(document)
  .ready(function() {

    /*
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
    */

    $('.ui.dropdown')
      .dropdown({
        on: 'hover'
      })
    ;

    /*
    $('.ui.form')
      .form(validationRules, {
        on: 'blur'
      })
    ;
    */

    $('.masthead .information')
      .transition('scale in', 1000)
    ;

    $('#menu-tabs .item')
        .tab('change tab', 'tab-info')
    ;

    $('#menu-tabs-info')
        .on('click', function() {
            loadInfoTab();
        })
    ;
    
    $('#menu-tabs-dockerversion')
        .on('click', function() {
            $.getJSON('/monitor/version')
                .done(function(data) {
                    document.getElementById('version-apiversion').innerHTML = data.ApiVersion;
                    document.getElementById('version-arch').innerHTML = data.Arch;
                    document.getElementById('version-gitcommit').innerHTML = data.GitCommit;
                    document.getElementById('version-goversion').innerHTML = data.GoVersion;
                    document.getElementById('version-kernelversion').innerHTML = data.KernelVersion;
                    document.getElementById('version-os').innerHTML = data.Os;
                    document.getElementById('version-version').innerHTML = data.Version;
                })
                .fail(function( jqxhr, textStatus, error ) {
                    var err = textStatus + ', ' + error;
                    console.log('Request Failed: ' + err );
                })
            ;
        })
    ;

    $('#menu-tabs-ping')
        .on('click', function() {
            $.getJSON("/monitor/ping")
                .done(function(data) {
                    $('#tab-ping #results').text('Ping response: ' + JSON.stringify(data));
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

    function loadInfoTab() {
        if (nodeId !== '') {
            console.log('Load info with nodeId: ' + nodeId);
        }
        $.getJSON("/monitor/info")
            .done(function(data) {
                var driverStatus = '';
                for (var i=0; i<data.DriverStatus.length; i++) {
                    driverStatus += '<br />';
                    for (var j=0; j<data.DriverStatus[i].length; j++) {
                        driverStatus += '&nbsp;&nbsp;&nbsp;&nbsp;' + data.DriverStatus[i][j] + ' : ';
                    }
                }
                var registryConfig = '';
                registryConfig += '<br />&nbsp;&nbsp;&nbsp;&nbsp;Index Configs: ' + 
                    JSON.stringify(data.RegistryConfig.IndexConfigs);
                registryConfig += '<br />&nbsp;&nbsp;&nbsp;&nbsp;Insecure Registry CIDRs: ' + 
                    JSON.stringify(data.RegistryConfig.InsecureRegistryCIDRs);
                document.getElementById('info-containers').innerHTML = data.Containers;
                document.getElementById('info-debug').innerHTML = data.Debug;
                document.getElementById('info-dockerrootdir').innerHTML = data.DockerRootDir;
                document.getElementById('info-driver').innerHTML = data.Driver;
                document.getElementById('info-driverstatus').innerHTML = driverStatus;
                document.getElementById('info-executiondriver').innerHTML = data.ExecutionDriver;
                document.getElementById('info-id').innerHTML = data.ID;
                document.getElementById('info-ipv4forwarding').innerHTML = data.IPv4Forwarding;
                document.getElementById('info-images').innerHTML = data.Images;
                document.getElementById('info-indexserveraddress').innerHTML = data.IndexServerAddress;
                document.getElementById('info-initpath').innerHTML = data.InitPath;
                document.getElementById('info-initsha1').innerHTML = data.InitSha1;
                document.getElementById('info-kernelversion').innerHTML = data.KernelVersion;
                document.getElementById('info-labels').innerHTML = data.Labels;
                document.getElementById('info-memtotal').innerHTML = data.MemTotal;
                document.getElementById('info-memorylimit').innerHTML = data.MemoryLimit;
                document.getElementById('info-ncpu').innerHTML = data.NCPU;
                document.getElementById('info-neventslistener').innerHTML = data.NEventsListener;
                document.getElementById('info-nfd').innerHTML = data.NFd;
                document.getElementById('info-ngoroutines').innerHTML = data.NGoroutines;
                document.getElementById('info-name').innerHTML = data.Name;
                document.getElementById('info-operatingsystem').innerHTML = data.OperatingSystem;
                document.getElementById('info-registryconfig').innerHTML = registryConfig;
                document.getElementById('info-swaplimit').innerHTML = data.SwapLimit;
                document.getElementById('info-systemtime').innerHTML = data.SystemTime;
            })
        .fail(function( jqxhr, textStatus, error ) {
            var err = textStatus + ", " + error;
            console.log( "Request Failed: " + err );
        })
        ;
    }

    var nodeId = document.getElementById('hidden-nodeid').value;

    loadInfoTab();

  })
;
