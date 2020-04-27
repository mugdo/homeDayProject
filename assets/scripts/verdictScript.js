var url = document.getElementById("url").innerText
console.log(url)

$(document).ready(function () {
    var doCheck = setInterval(function () {
        var counter = 0;
        $.getJSON(url, function (result) {
            $('#res').text(result.status);
            counter++;
            console.log(counter)
            if (counter > 15 || result.status == "Accepted" || result.status == "Worng Answer") {
                clearInterval(doCheck);
            }
        });
    }, 2000);  //Delay here = .5 seconds
});