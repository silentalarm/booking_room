var today = new Date();
var date = setZero(today.getDate()) + '.' + (setZero(today.getMonth() + 1)) + '.' + today.getFullYear()

// одиночное выделеление
var lineSet = new Set();

//защита от ввода букв в поле для количества человек
function validateNumber(event) {
    var ASCIICode = (event.which) ? event.which : event.keyCode
    if (ASCIICode > 31 && (ASCIICode < 48 || ASCIICode > 57))
        return false;
    return true;
}
//для даты и времени (все что меньше 10 показывалось без 0 вначале)
function setZero(someVar) {
    if (someVar < 10)
        return (someVar = "0" + someVar);
    else
        return (someVar);
}

// проверка на наличие чисел больше 0
function isMoreThenZero(num) {
    if (num > '0' && num.length > 0)
        return true
    else
        return false
}

$('#room_select').on('change', function () {
    var floor = $(this).val();//этаж
    //var dat = date //today
    var date = $('#dp').datepicker('getFormattedDate');
    if (floor) {
        window.location = "?table="+floor + "&date=" + date;
    }
    return false;

});



// переделать поля на заполненность формы
clubName = $('#clubName').val()

if (clubName != null) {
    function checkHandler() {
        if ($('#statusAuth').val() == "true" && $('#clubName').val().trim().length > 0 && lineSet.size > 0 && isMoreThenZero($('#peopleNumber').val().trim())) {
            $('#btn_reservation').prop('disabled', false);
        } else {
            $('#btn_reservation').prop('disabled', true);
        }
    }
    $(document).on('click', '.tr_size', function () {
        if ($(this).hasClass('bg-select')) {
            $('.tr').text($(this).index());
            lineSet.delete($(this).index());
            $(this).removeClass('bg-select');
        } else {
            if ($(this).hasClass('bg-booked')) {
                return // не учитываем забронированные
            }
            $('.tr').text($(this).index());
            lineSet.add($(this).index());
            $(this).addClass('bg-select');
        }
        // пишем в input lines наши одиночно выделенные строки (то есть часы)
        $('#lines').val(Array.from(lineSet));
        checkHandler();
    });

// мультивыделение
    $(document).on('mousedown', '.tr_size', function () {
        if ($(this).hasClass('bg-select')) {
            $('.tr_size').hover(function () {
                $('.tr').text($(this).index());
                lineSet.delete($(this).index());
                $(this).removeClass('bg-select');
            });
        } else {
            $('.tr_size').hover(function () {
                if ($(this).hasClass('bg-booked')) {
                    return // не учитываем забронированные
                }
                $('.tr').text($(this).index());
                lineSet.add($(this).index());
                $(this).addClass('bg-select');
            });
        }
    });

    $('.tr_size').mouseup(function () {
        // пишем в input lines наши мультивыделенные строки (то есть часы)
        checkHandler();
        $('#lines').val(Array.from(lineSet));
        $('.tr_size').off('mouseenter mouseleave');
    });

    $(".tr_size").on("mousedown", function (e) {
        e.preventDefault();
        $(this).addClass("pointer");
    }).on("mouseup", function () {
        $(this).removeClass("pointer");
    });
}

//удаление резерва
$(document).on('click', '.unreserve', function () {
    var date = $('#dp').datepicker('getFormattedDate');
    var floor = $('#room_select').val();
    var time = $((this).parentElement.parentElement).index()
    var club = $('#clubName').val();

    window.location = "/delreserve?table="+floor + "&date=" + date + "&deltime=" + time + "&clubname=" + club;
});

//при нажатии показывает описание клуба
$(".row_club").on("click",
    function() {
        var accordionRow = $(this).next(".slider");
        if (!accordionRow.is(":visible")) {
            accordionRow.show().find(".slider-content").slideDown();
        } else {
            accordionRow.find(".slider-content").slideUp(function() {
                if (!$(this).is(':visible')) {
                    accordionRow.hide();
                }
            });
        }
});

$(function () {
    $('#cp2').colorpicker();
});



$(document).ready(function () {

    $('.notes').on('click',function (e){
        alert("ok");
        $.ajax({
            type:'GET',
            url :'localhost:8185/?table=floor_3&date=17.12.2020',
            dataType: 'html',
            success: function(data) {
                console.log('success',data);
                $('#interactive').html(data);

            },
            error: function(jqXHR,textStatus,errorThrown ){
                alert('Exception:'+errorThrown );
            }
        });
        e.preventDefault();
    });

 $('#room_select').val($('#selectedFloor').val())

    $('[data-toggle="tooltip"]').tooltip();

    setInterval(function () {
        var today = new Date();
        var date = setZero(today.getDate()) + '.' + (setZero(today.getMonth() + 1)) + '.' + today.getFullYear();
        var dateTime = date;

        $('#date').val(dateTime);
    }, 1000);

    date1 = new Date(2011, 1, 28)


    $('#dp').datepicker({
        format: "dd.mm.yyyy",
        startDate: "today",
        endDate: "+30d",
        maxViewMode: 0,
        language: "ru",
        weekStart: 1,
        daysOfWeekHighlighted: "0,6",
        todayHighlight: true,
        multidate: false,
        keyboardNavigation: false,
        forceParse: false,
        //toggleActive: true,
        assumeNearbyYear: true,
        startDate: "today",
        setDate: date1 ,
    }).on('changeDate', showTestDate).datepicker("update", $('#selectedDate').val());

    function showTestDate(){
        var date = $('#dp').datepicker('getFormattedDate');
        var floor = $('#room_select').val();

        $("#showDate").val(date);
        $('#dp').datepicker('select', $('#selectedDate').val());
        window.location = "?table="+floor + "&date=" + date;
    }

$('.datepicker-switch').setAttribute('disabled', true);

})


$('input[type="file"]').change(function(e){
    var fileName = e.target.files[0].name;
    $('.custom-file-label').html(fileName);
    onFileSelected(e)
});

function onFileSelected(event) {
    var selectedFile = event.target.files[0];
    var reader = new FileReader();

    var imgtag = document.getElementById("clubPic");
    imgtag.title = selectedFile.name;
    $('#imageInputButton').prop('disabled', false);
    reader.onload = function(event) {
        imgtag.src = event.target.result;
    };

    reader.readAsDataURL(selectedFile);
}

$(document).ready(function(){
    $('.square').mouseover(function (e) {
        child = $(this).find('.trgt')
        child.show();

        var all = $(window).width();
        var left = $(this).offset().left + 150;
        var width = $(this).outerWidth(true);
        var offset = all - (left + width);
        console.log(offset)
        $(this).removeClass('right-tl');
        $(this).removeClass('left-tl');
        if (offset < 0)
        {
            $(this).addClass("right-tl");
        }else{
            $(this).addClass("left-tl");
        }
    });

    $('.square').mouseout(function (e) {
        child = $(this).find('.trgt')
        child.hide();
    });
});

$(document).ready(function(){
    $('.square').on('click',function (e) {
        child = $(this).find('.trgt')
        child.show();
    });
});
function changeColor(value){
    document.getElementById("color").value = value;
    document.getElementById("colorInp").value = value;
}

function changeLineColor(id ,value){
    document.getElementById(id).style.backgroundColor = value;
}

$("#colorGroup").change(function(){
    var color = document.getElementById("colorGroup").value;
    document.getElementById("colorGroup").style.backgroundColor = color;
});


$('#clubName').on('change', function () {
    var clubName = $(this).val();
    console.log(clubName)


    $.ajax({
        url: "/getgroups",
        method: "POST",
        contentType: 'application/json; charset=utf-8',
        data: JSON.stringify({ Name: clubName}),
        dataType: 'json',
        success: function(r) {
            var myDiv = document.getElementById("groupName");
            var array = r.List;
;
            for (var i = 0; i < array.length; i++) {
                var option = document.createElement("option");
                option.value = array[i];
                option.text = array[i];
                myDiv.appendChild(option);
            }
        }
    });
});

$(document).ready(function(){

});