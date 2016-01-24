$(document).ready(function () {/* activate scrollspy menu */
    $('body').scrollspy({
        target: '#navbar-collapsible',
        offset: 50
    });

    /* smooth scrolling sections */
    $('a[href*=#]:not([href=#])').click(function () {
        if (location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '') && location.hostname == this.hostname) {
            var target = $(this.hash);
            target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
            if (target.length) {
                $('html,body').animate({
                    scrollTop: target.offset().top - 50
                }, 1000);
                return false;
            }
        }
    });

    var username = getCookie("username")
    if (username) {
        $("#login-info").html(
            "<a href='/dashboard'>" +
            "<img src='"+getCookie("useravatar")+"'>&nbsp" +
            username + "</a>"
        );
    }

    $('.fav').click(function () {

        var icon = $('.fav').find('i');
        var action = icon.hasClass("fa-heart-o") ? "add" : "remove";
        var title = $(this).parent().prev().text().trim();
        $.ajax("/fav/" + action, {
            method: "POST",
            data: {package: title},
            dataType: "json"
        });

        if (action == "add") {
            icon.removeClass("fa-heart-o").addClass("fa-heart");
        } else {
            icon.removeClass("fa-heart").addClass("fa-heart-o");
        }

        if (typeof $(this).parents('.package') !== 'undefined' ) {
            $(this).parents('.package').remove();
        }
    });

    $('.package').each(function () {
            var data = $(this).data("bench").split(",").map(function (e) {return parseInt(e)});
            var title = $(this).find(".box-title").text();
            $(this).find('.chart-container').highcharts({
                chart: {
                    type: 'area'
                },
                title: {
                    text: ''
                },
                xAxis: {
                    categories: ['', '', '', '', '', '', ''],
                    tickmarkPlacement: 'on',
                    title: {
                        enabled: false
                    }
                },
                yAxis: {
                    title: {
                        text: ''
                    }
                },
                plotOptions: {
                    area: {
                        stacking: 'normal',
                        lineColor: '#666666',
                        lineWidth: 1,
                        marker: {
                            lineWidth: 1,
                            lineColor: '#666666'
                        }
                    }
                },
                series: [{
                    name: title,
                    data: data
                }]
            });
        });
});



function getCookie(name) {
    var matches = document.cookie.match(new RegExp(
        "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
    ));
    return matches ? decodeURIComponent(matches[1]) : undefined;
}
