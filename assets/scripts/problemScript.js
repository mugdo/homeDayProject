console.log("Script linked properly")

function processPList(getPList) {
    pList = getPList;
}

var pageUL = $('#pageUL')
$(window).on("load", function () {
    var counter=1;
    for (i = 0; i < pList.length; i+=20) {
        pageUL.append(`<li><a id="page`+counter+`" href="#" onclick="showProblem(`+counter+`)">`+counter+`</a></li>`);
        counter++;
    }
    showProblem(1);
});
function showProblem(pageNum) {
    removeRow();

    var start = 20*(pageNum-1);
    var finish = start + 20;

    for (i = start; i < Math.min(finish, pList.length); i++) {
        var link = "/problemView/" + pList[i].originOJ + "-" + pList[i].originProb;

        $('#problemTable').append(`<tr class="problemRow">
            <td>`+ pList[i].originOJ + `</td>
            <td>`+ pList[i].originProb + `</td>
            <td><a href="`+ link + `">` + pList[i].title + `</a></td>
        </tr>`);
    }
}
function removeRow() {
    var rowSize = document.getElementById("problemTable").rows.length;

    for (i = 0; i < rowSize - 2; i++) {
        $('.problemRow').remove();
    }
}

{
    {/* $(document).ready(function() {
    }); */}
}