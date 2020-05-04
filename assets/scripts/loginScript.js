console.log("Login Script linked properly")

var username = $('#username')
var password = $('#password')

var errUsername = $('#errUsername')
var errPassword = $('#errPassword')

username.keyup(function () {
    if (username.val().length == 0) {
        errUsername.text("")
    } else {
        $.ajax({
            url: "/check/username=" + username.val().trim(),
            method: "GET",
            success: function (data) {
                if (data == "true") {
                    errUsername.text("")
                } else {
                    errUsername.text("Username not found.")
                }
            },
        });
    }
});
password.keyup(function () {
    errPassword.text("")
});

$('.form-group').keyup(function () {
    var checkSubmit = setInterval(function () {
        if (errUsername.text() === "") {
            $('#submit').removeAttr('disabled');
            $('#submit').removeAttr('Title');
        } else {
            $('#submit').attr('disabled', 'disabled');
            $('#submit').attr('Title', 'Please fix the error.');
        }
        clearInterval(checkSubmit)
    }, 100);
});

{
    {/* $(document).ready(function() {
    }); */}
}