console.log("Script linked properly")
var OJReal = document.getElementById('OJReal').innerText
var OJBox = $('#OJBox')
var languageBox = $('#submit-language')

function addSelected(need) {
    OJBox.find('option').each(function () {
        if ($(this).text() == need) {
            $(this).attr("selected", "selected");
            return false;
        }
    });
};
function removeSelected() {
    OJBox.find('option').each(function () {
        $(this).removeAttr('selected');
        return false;
    });
};
function removeLanguage() {
    languageBox.find('option').each(function () {
        $(this).remove();
    });
};
function addLanguage(OJ) {
    $.ajax({
        url: "/lang=" + OJ,
        method: "GET",
        success: function (data) {
            console.log("success")

            $.each(data, function (key, value) {
                languageBox.append($("<option></option>")
                    .attr("value", value.LangValue)
                    .text(value.LangName));
            });
        },
    });
};

addSelected(OJReal); //adding selected tag
addLanguage(OJReal);

OJBox.on('change', function () {
    removeSelected();
    addSelected(this.value);

    var selectedOJ = $('select[name="OJ"]').val()
    removeLanguage();
    addLanguage(selectedOJ);
});

{
    {/* $(document).ready(function() {
    }); */}
}