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

    /*
    $('.menu .item')
      .tab()
    ;
    */
    
    $('.menu .item')
        .tab({
            cache: false,
            // faking api request
            apiSettings: {
                loadingDuration : 3000
                , 
                onLoad: function(x, y, response) {
                    console.log('ok');
                    var response = {
                        results : []
                    }
                    ;
                    return response;
                }
            // ,
            //     mockResponse    : function(settings) {
            //         var response = {
            //             list : 'AJAX Tab One',
            //             inspect : 'AJAX Tab Two',
            //             search : 'AJAX Tab Three'
            //         };
            //         return response[settings.urlData.tab];
            //     }
            },
            context : 'parent',
            auto    : true,
            path    : '/images/'
        })
    ;
    
    $.fn.api.settings.successTest = function(response) {
        // if(response && response.success) {
        //     return response.success;
        // }
        // return false;
        return true;
    };

// Function to create a table as a child of el.
// // data must be an array of arrays (outer array is rows).
    function tableCreate(el, data)
    {
        var tbl  = document.createElement("table");
        tbl.style.width  = "70%";

        for (var i = 0; i < data.length; ++i)
        {
            var tr = tbl.insertRow();
            for(var j = 0; j < data[i].length; ++j)
            {
                var td = tr.insertCell();
                td.appendChild(document.createTextNode(data[i][j].toString()));
            }
        }
        el.appendChild(tbl);
    }

    /*
    $.post("/whatever", { somedata: "test" }, null, "json")
        .done(function(data) {
            rows = [];
            for (var i = 0; i < data.Results.length; ++i)
        {
            cells = [];
            cells.push(data.Results[i].A);
            cells.push(data.Results[i].B);
            rows.push(cells);
        }
        tableCreate($("#results")[0], rows);
        });
    */

    /*
    $('.infinite.example .demo.segment')
      .visibility({
        once: false,
              // update size when new content loads
              //     observeChanges: true,
              //         // load content on bottom edge visible
              //             onBottomVisible: function() {
              //                   // loads a max of 5 times
              //                         window.loadFakeContent();
              //                             }
      })
      ;
      */

    setInterval(changeSides, 3000);

  })
;
