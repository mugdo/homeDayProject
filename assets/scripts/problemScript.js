console.log("Script linked properly")

let shiftX = 32, previousPage = 5, currPage = 1;
function processPList(getPList) {
    console.log(getPList)
    pList = getPList;
    totalPages = Math.ceil(pList.length / 20); //20 problem per page

    displayingPageNum = Math.min(9, totalPages) //by default 9 pageNum displayed on pagination
    rightBound = (totalPages - displayingPageNum) * -(shiftX);
}

var pageUL = $('#pageUL')
$(window).on("load", function () {
    //page number creation
    for (i = 1; i <= totalPages; i++) {
        pageUL.append(`<li><a id="pageBox` + i + `" href="#" onclick="showProblem(` + i + `)">` + i + `</a></li>`);
    }
    //width fixing of pagination
    pagerWidth = displayingPageNum * 32
    $("#pagination").css("width", pagerWidth + "px");
    $("#pager").css("width", (pagerWidth + 140) + "px"); //extra 140px for 'pre/next' button

    //onload-showing problem of  page 1
    showProblem(1);
});

function showProblem(active) {
    currPage = active;
    //removing current existing rows
    var rowSize = document.getElementById("problemTable").rows.length;
    for (i = 0; i < rowSize - 2; i++) {
        $('.problemRow').remove();
    }

    //adding new rows
    var start = 20 * (active - 1);
    var finish = start + 20;

    for (i = start; i < Math.min(finish, pList.length); i++) {
        var link = "/problemView/" + pList[i].originOJ + "-" + pList[i].originProb;

        $('#problemTable').append(`<tr class="problemRow">
            <td>`+ pList[i].originOJ + `</td>
            <td>`+ pList[i].originProb + `</td>
            <td><a href="`+ link + `">` + pList[i].title + `</a></td>
        </tr>`);
    }
    updateActiveClass(active);
    shifting(previousPage, active);
}
function updateActiveClass(currPage) {
    for (i = 1; i <= totalPages; i++) {
        id = "#pageBox" + i;

        if ($(id).attr("class")) {
            $(id).removeClass("activePage disableClick")
        }

        if ($(id).text() == currPage) {
            $(id).addClass("activePage disableClick")
        }
    }

    //previous button
    if (currPage == 1) {
        $('#pre').addClass("disableClick")
    } else {
        $('#pre').removeClass("disableClick")
    }
    //next button
    if (currPage == totalPages) {
        $('#nxt').addClass("disableClick")
    } else {
        $('#nxt').removeClass("disableClick")
    }
}
function shifting(previousPage, currPage) {
    if (currPage == 1 || currPage == 2 || currPage == 3 || currPage == 4 || currPage == 5) {
        $("#pageUL").css({ "transform": "translateX(" + -0 + "px)", "transition": "transform 1s ease" });
    } else if (currPage == totalPages || currPage == (totalPages - 1) || currPage == (totalPages - 2) || currPage == (totalPages - 3) || currPage == (totalPages - 4)) {
        $("#pageUL").css({ "transform": "translateX(" + rightBound + "px)", "transition": "transform 1s ease" });
    } else {
        diff = currPage - 5;
        $("#pageUL").css({ "transform": "translateX(" + -(diff * 32) + "px)", "transition": "transform 1s ease" });
    }
    previousPage = currPage;
}
$('#pre').click(function () {
    showProblem(currPage - 1);
});
$('#nxt').click(function () {
    showProblem(currPage + 1);
});
// $(document).ready(function () {
// });