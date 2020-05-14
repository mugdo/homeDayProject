var verdict = $('#verdict')
var time = $('#time')
var memory = $('#memory')
var submitTime = $('#submitTime')

var subID = document.getElementById("subID").innerText;
var url = "/verdict/subID="+subID;
console.log(url)

$(document).ready(function () {
    var counter = 0;
    var doCheck = setInterval(function () {
        $.getJSON(url, function (result) {
            counter++;
            if (result.status == "Accepted" || result.status == "Wrong Answer" || result.status == "Compilation Error" || result.status == "Time Limit Exceeded" || result.status == "Memory Limit Exceeded") {
                
                verdict.text(result.status)
                time.text(result.runtime)
                memory.text(result.memory)
                submitTime.text(result.submitTime)

                clearInterval(doCheck);
            } else if (counter > 30) {
                //location.reload();
                clearInterval(doCheck);
            }
        });
    }, 2000);  //Delay here = 2 seconds
});