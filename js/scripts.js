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
    var dat = date //today
    if (floor) {
        window.location = "?table="+floor + "&date=" + dat;
    }
    return false;

});


// переделать поля на заполненность формы
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


//удаление резерва
$(document).on('click', '.unreserve', function () {
    alert($((this).parentElement.parentElement).index());
});


//$('[data-toggle="datepicker"]').datepicker({
//    format: 'dd.mm.yyyy',
//    weekStart: 1,
//    startDate: date,
//    autoPick: true,
//    autoShow: true,
//    inline: true,
//    container: '.datepicker-inline',
//    months: ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'],
//    monthsShort: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'],
//    days: ['воскресенье', 'понедельник', 'вторник', 'среда', 'четверг', 'пятница', 'суббота'],
//    daysShort: ['вск', 'пнд', 'втр', 'срд', 'чтв', 'птн', 'сбт'],
//    daysMin: ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'],
//});
//
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



 // заполнение select
 //$("select option[value=" + val + "]").attr('selected', 'true').text(text);

 $('#room_select').val($('#selectedFloor').val())

// tooltip
    $('[data-toggle="tooltip"]').tooltip();

    //login logout
    var status = $('#statusAuth').val()
    console.log("Authentication is " + status)
    if (status == "false") {
        $('#logout').attr("hidden", true)
        $('#login').attr("hidden", false)
    } else {
        $('#login').attr("hidden", true)
        $('#logout').attr("hidden", false)
    }

    /* Локализация datepicker */
//$.datepicker.regional['ru'] = {
//    closeText: 'Закрыть',
//    prevText: 'Предыдущий',
//    nextText: 'Следующий',
//    currentText: 'Сегодня',
//    monthNames: ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'],
//    monthNamesShort: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'],
//    dayNames: ['воскресенье', 'понедельник', 'вторник', 'среда', 'четверг', 'пятница', 'суббота'],
//    dayNamesShort: ['вск', 'пнд', 'втр', 'срд', 'чтв', 'птн', 'сбт'],
//    dayNamesMin: ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'],
//    weekHeader: 'Не',
//    dateFormat: 'dd.mm.yy',
//    firstDay: 1,
//    isRTL: false,
//    showMonthAfterYear: false,
//    yearSuffix: ''
//};
//$.datepicker.setDefaults($.datepicker.regional['ru']);

//('#datepicker').datepicker([options])
//$.fn.datepicker.noConflict();


//$(function(){
//    $("#datepicker").datepicker({
//        onSelect: function(date){
//            $('#datepicker_value').val(date)
//        }
//    });
//    $("#datepicker").datepicker("setDate", $('#datepicker_value').val());
//
//});


////liveTime
//     setInterval(function () {
//         // Just move your date creation inside the interval function
//         var today = new Date();
//         var date = setZero(today.getDate()) + '.' + (setZero(today.getMonth() + 1)) + '.' + today.getFullYear();
//         var time = setZero(today.getHours()) + ":" + setZero(today.getMinutes());// + ":" + today.getSeconds();
//         var dateTime = date + " " + time; // Add the time to the date string
//
//         document.getElementById('clock').innerHTML = dateTime;
//     }, 1000);

    // пишем в input дату
    setInterval(function () {
        // Just move your date creation inside the interval function
        var today = new Date();
        var date = setZero(today.getDate()) + '.' + (setZero(today.getMonth() + 1)) + '.' + today.getFullYear();
        //var time = setZero(today.getHours()) + ":" + setZero(today.getMinutes());// + ":" + today.getSeconds();
        var dateTime = date; // Add the time to the date string

        $('#date').val(dateTime);
    }, 1000);

//
//$('#datepicker').datepicker({
//inline: true,
//container: '.datepicker-inline',
//    format: "dd.mm.yyyy",
//    language: "ru",
//    weekStart: 1,
//    daysOfWeekHighlighted: "0,6",
//    todayHighlight: true,
//    //autoclose: true,
//});
//$('#datepicker').datepicker("setDate", new Date());
//
   // $('#sandbox-container div').datepicker({
   // });
//
   // $('#datepicker').datepicker({
   //     weekStart: 1,
   //     daysOfWeekHighlighted: "6,0",
   //     autoclose: true,
   //     todayHighlight: true,
   // });
   // $('#datepicker').datepicker("setDate", new Date());
//
//

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

    //$('#datepicker').datepicker("setDate", new Date());

    function showTestDate(){
        var value = $('#dp').datepicker('getFormattedDate');
        $("#showDate").val(value);
        $('#dp').datepicker('select', $('#selectedDate').val());
        //alert($('#selectedDate').val())

        //$('#dp').datepicker("update", date1);
        //alert($('#selectedDate').val())


        var floor = $('#room_select').val();//этаж
        window.location = "?table="+floor + "&date=" + value;

    }

$('.datepicker-switch').setAttribute('disabled', true);

})
