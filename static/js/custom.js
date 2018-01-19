$(function () {
    $('.navbar-right .dropdown-menu .body .menu').slimscroll({
        height: '254px',
        color: 'rgba(0,0,0,0.5)',
        size: '4px',
        alwaysVisible: false,
        borderRadius: '0',
        railBorderRadius: '0'
    });

    $("a[data-confirm]").on('click', function () {
        var $a = $(this);
        swal({
            title: "Are you sure?",
            text: "You will not be able to undo this operation!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#c9302c",
            confirmButtonText: "Yes, do it!",
            closeOnConfirm: false
        }, function () {
            window.location = $a.attr("href");
        });
        return false;
    });

    $('form').validate({
        highlight: function (input) {
            $(input).parents('.form-line').addClass('error');
        },
        unhighlight: function (input) {
            $(input).parents('.form-line').removeClass('error');
        },
        errorPlacement: function (error, element) {
            $(element).parents('.form-group').append(error);
        }
    });
});