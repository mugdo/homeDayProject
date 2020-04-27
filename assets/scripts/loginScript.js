console.log("Login Script linked properly")

var username = $('#username')
var errUsername = $('#errUsername')

username.keyup(function () {
    console.log("In kup")
    $.ajax({
        url: "/check/username=" + username.val().trim(),
        method: "GET",
        success: function (data) {
            if (data == "true") {
                console.log("true")
                errUsername.text("")
            } else {
                console.log("false")
                errUsername.text("Username not found.")
            }
        },
    });
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