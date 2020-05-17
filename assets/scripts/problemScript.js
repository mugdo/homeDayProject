console.log("Script linked properly")

function processPList(getPList) {
    console.log(getPList)
    pList = getPList;
    totalPages = Math.ceil(pList.length / 20);
}

var pageUL = $('#pageUL')
$(window).on("load", function () {
    showProblem(1);
});

function showProblem(active) {
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

    resetPageNum(active);
}
function resetPageNum(active) {
    removeExistingPageNum();
    addingCurrentPageNum(active);

    updateActiveClass(active);
    updateDisableClass(active);
}
function removeExistingPageNum() {
    $('#pageBoxPre').parent().remove();
    $('#pageBox1').parent().remove();
    $('#pageBox2').parent().remove();
    $('#pageBox3').parent().remove();
    $('#pageBox4').parent().remove();
    $('#pageBox5').parent().remove();
    $('#pageBoxNext').parent().remove();
}
function addingCurrentPageNum(active) {
    //prevoius button
    pageUL.append(`<li><a id="pageBoxPre" href="#" onclick="showProblem(` + (active - 1) + `)">Previous</a></li>`);

    //middle 5 button
    if (totalPages <= 5) {
        for (i = 1; i <= totalPages; i++) {
            pageUL.append(`<li><a id="pageBox` + i + `" href="#" onclick="showProblem(` + i + `)">` + i + `</a></li>`);
        }
    } else {
        let pageCutLow = active - 2;
        let pageCutHigh = active + 2;

        if (pageCutLow < 1) {
            id = 1;
            for (i = 1; i <= 5; i++) {
                pageUL.append(`<li><a id="pageBox` + id + `" href="#" onclick="showProblem(` + i + `)">` + i + `</a></li>`);
                id++;
            }
        } else if (pageCutHigh > totalPages) {
            id = 1;
            for (i = (totalPages - 4); i <= totalPages; i++) {
                pageUL.append(`<li><a id="pageBox` + id + `" href="#" onclick="showProblem(` + i + `)">` + i + `</a></li>`);
                id++;
            }
        } else {
            id = 1;
            for (i = pageCutLow; i <= pageCutHigh; i++) {
                pageUL.append(`<li><a id="pageBox` + id + `" href="#" onclick="showProblem(` + i + `)">` + i + `</a></li>`);
                id++;
            }
        }
    }

    //next button
    pageUL.append(`<li><a id="pageBoxNext" href="#" onclick="showProblem(` + (active + 1) + `)">Next</a></li>`);
}
function updateActiveClass(active) {
    for (i = 1; i <= 5; i++) {
        id = "#pageBox" + i;
        $(id).removeClass("activePage")

        if ($(id).text() == active) {
            $(id).addClass("activePage disablePage")
        }
    }
}
function updateDisableClass(active) {
    if (active == 1) {
        $('#pageBoxPre').addClass("disablePage")
    } else {
        $('#pageBoxPre').removeClass("disablePage")
    }
    if (active == totalPages) {
        $('#pageBoxNext').addClass("disablePage")
    } else {
        $('#pageBoxNext').removeClass("disablePage")
    }
}
// $(document).ready(function () {
// });
// $('#pre').click(function(){
//     $("#pageUL").css({ "transform":"translateX(50px)","transition": "transform 1s ease" });
// });
// $('#next').click(function(){
//     $("#pageUL").css({ "transform":"translateX(50px)","transition": "transform 1s ease" });
// });
