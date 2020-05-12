var url = document.getElementById("url").innerText
console.log(url)

$(document).ready(function () {
    var counter = 0;
    var doCheck = setInterval(function () {
        $.getJSON(url, function (result) {
            $('#res').text(result.status);
            counter++;
            if (counter > 15 || result.status == "Accepted" || result.status == "Wrong Answer" || result.status == "Compilation Error" || result.status == "Time Limit Exceeded") {
                console.log(counter+result.status)
                clearInterval(doCheck);
            }
        });
    }, 2000);  //Delay here = 2 seconds
});