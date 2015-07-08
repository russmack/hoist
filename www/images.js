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
                loadingDuration : 3000,
                mockResponse    : function(settings) {
                    var response = {
                        list : 'AJAX Tab One',
                        inspect : 'AJAX Tab Two',
                        search : 'AJAX Tab Three'
                    };
                    return response[settings.urlData.tab];
                }
            },
            context : 'parent',
            auto    : true,
            path    : '/images/'
        })
    ;
    


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
