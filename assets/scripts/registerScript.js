console.log("Script linked properly")

var username = $('#username')
var email = $('#email')
var password = $('#password')
var confirmPassword = $('#confirmPassword')

var errUsername = $('#errUsername')
var errEmail = $('#errEmail')
var errPassword = $('#errPassword')

username.keyup(function () {
    $.ajax({
        url: "/check/username=" + username.val().trim(),
        method: "GET",
        success: function (data) {
            if (data == "true") {
                errUsername.text("Username already taken. Choose another one.")
            } else {
                errUsername.text("")
            }
        },
    });
});

email.keyup(function () {
    $.ajax({
        url: "/check/email=" + email.val().trim(),
        method: "GET",
        success: function (data) {
            if (data == "true") {
                errEmail.text("Email already registered. Choose another one.")
            } else {
                errEmail.text("")
            }
        },
    });
});

password.keyup(function () {
    var pass = password.val()
    var confirmPass = confirmPassword.val()

    if (confirmPass.length > 0 || (confirmPass.length == 0 && errPassword.text() != "")) {
        if (pass !== confirmPass) {
            errPassword.text("Password mismatched. Put cautiously")
        } else {
            errPassword.text("")
        }
    }
});
confirmPassword.keyup(function () {
    var pass = password.val()
    var confirmPass = confirmPassword.val()

    if (pass !== confirmPass) {
        errPassword.text("Password mismatched. Put cautiously")
    } else {
        errPassword.text("")
    }
});

$('.form-group').keyup(function () {
    var checkSubmit = setInterval(function () {
        if (errUsername.text() === "" && errEmail.text() === "" && errPassword.text() === "") {
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