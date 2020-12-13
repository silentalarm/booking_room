// $.fn.datepicker.noConflict();
// $.fn.datepicker.languages['ru-RU'] = {
//     format: 'dd.mm.YYYY',
//     days: ['Воскресенье', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота'],
//     daysShort: ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'],
//     daysMin: ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'],
//     months: ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'],
//     monthsShort: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'],
//     weekStart: 1,
//     startView: 0,
//     yearFirst: false,
//     yearSuffix: ''
// };


//защита от ввода букв в поле для количества человек
// function validateNumber(event) {
//     var ASCIICode = (event.which) ? event.which : event.keyCode
//     if (ASCIICode > 31 && (ASCIICode < 48 || ASCIICode > 57))
//         return false;
//     return true;
// }


//для даты и времени (все что меньше 10 показывалось без 0 вначале)
function setZero(someVar) {
    if (someVar < 10)
        return (someVar = "0" + someVar);
    else
        return (someVar);
}

//
//Highlighting rows
//(function ($) {
//    $.fn.checkboxTable = function () {
//        target = this;
//        //click on checkbox.
//        $(target).on('click', 'tbody :checkbox', function () {
//            var parents = $(this).parents('table');
//            if ($(this).is(':checked')) {
//                $(this).parents('tr').addClass('bg-select');
//            } else {
//                $(this).parents('tr').removeClass('bg-select');
//                if ($(parents).find('tbody :checkbox:checked').length == 0) {
//                    $(this).closest('tr').find("input[type=text], text").val("");
//                    $(this).closest('tr').find("input[type=text], text").first().val("Disabled");
//                }
//            }
//        });
//    };
//})(jQuery);

//

function isMoreThenZero(num) {
    if (num > '0' && num.length > 0)
        return true
    else
        return false
}

function checkBoxHandler() {
    let index = 0;
    for (let i = 0; i < 24; i++) {
        if ($(`#checkBox${i}`).prop('checked') && $(`#clubName${i}`).val().trim().length > 0 && isMoreThenZero($(`#peopleNumber${i}`).val().trim())) {
            $('#btn_reservation').prop('disabled', false);
            index = i;
            break
        } else {
            index = 0;
            $('#btn_reservation').prop('disabled', true);
        }
    }

    //$(".triggerkkk").on("click", function() {
    //    $("#secondSelect").val($("#firstSelect").val());
    //})

    //for (let i = 0; i < 24; i++) {
    //    if ($(`#checkBox${i}`).prop('checked')) {
    //        $(`#clubName${i}`).val($(`#clubName${index}`).val())
    //        $(`#peopleNumber${i}`).val($(`#peopleNumber${index}`).val())
    //    }
    //}

}

//очистка inputs если чекбокс убран
// jQuery(function ($) {
//     $(document).on('change', '.check_box__status', function () {
//         //момент отжатия чекбокса
//         if (!this.checked) {
//             // console.log("Изменение")
//             $(this).closest('tr').find("input[type=text], text").val("");
//             $(this).closest('tr').find("input[type=text], text").first().val("Disabled");
//             $(this).parents('tr').removeClass('bg-select');
//         }
//         //момент нажатия на чекбокс
//         if (this.checked) {
//             $(this).parents('tr').addClass('bg-select');
//             //ClubName
//             let index = 0;
//             for (let i = 0; i < 24; i++) {
//                 if ($(`#checkBox${i}`).prop('checked')) {
//                     if ($(`#clubName${i}`).val() != '') {
//                         index = i;
//                         // console.log(index)
//                     }
//                 }
//             }
//             for (let i = 0; i < 24; i++) {
//                 if ($(`#checkBox${i}`).prop('checked')) {
//                     if ($(`#clubName${i}`).val() == '') {
//                         $(`#clubName${i}`).val($(`#clubName${index}`).val())
//                     }
//                 }
//             }
//             //Qantity
//             index = 0;
//             for (let i = 0; i < 24; i++) {
//                 if ($(`#checkBox${i}`).prop('checked')) {
//                     if ($(`#peopleNumber${i}`).val() != '') {
//                         index = i;
//                         // console.log(index)
//                     }
//                 }
//             }
//             for (let i = 0; i < 24; i++) {
//                 if ($(`#checkBox${i}`).prop('checked')) {
//                     if ($(`#peopleNumber${i}`).val() == '') {
//                         $(`#peopleNumber${i}`).val($(`#peopleNumber${index}`).val())
//                     }
//                 }
//             }
//         }
//     })
// });

// $('.cell').mousedown(function() {
//     $('#numa').text(parseInt($(this).parent('.row').index()));
//     a = parseInt($(this).parent('.row').index());
// });

// одиночное выделеление
let lineSet = new Set();
$(document).on('click', '.tr_size', function () {
    if ($(this).hasClass('bg-select')) {
        $('.tr').text($(this).index());
        lineSet.delete($(this).index());
        $(this).removeClass('bg-select');
    } else {
        $('.tr').text($(this).index());
        lineSet.add($(this).index());
        $(this).addClass('bg-select');
    }
    console.log(lineSet)
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
            $('.tr').text($(this).index());
            lineSet.add($(this).index());
            $(this).addClass('bg-select');
        });
    }
});
$('.tr_size').mouseup(function () {
    $('.tr_size').off('mouseenter mouseleave');
});


//if($(this).hasClass('test1')){
//    alert('У этого блока есть класс test1');
//}else{
//    alert('У этого блока нет класса test1');
//}

//if($(this).prop("checked") == true){
//    console.log("Checkbox is checked.");
//}
//else if($(this).prop("checked") == false){
//    console.log("Checkbox is unchecked.");
//}

//$(this).closest('tr').find("input[type=text], text").val("");
//$(this).closest('tr').find("input[type=text], text").first().val("Disabled");
//$(this).parents('tr').removeClass('bg-select');
//
// function changeBgColor(){
//     console.log("внешний клик")
//     var tr = document.getElementsByClassName('tr_size')
//     $('tr').click( function() {
//         console.log("внутренний клик")
//         $(this).addClass('bg-select')
//     } );
//}

//$( function() {
//    $( "#datepicker" ).datepicker();
//} );
var today = new Date();
var date = setZero(today.getDate()) + '.' + (setZero(today.getMonth() + 1)) + '.' + today.getFullYear()

//$('[data-toggle="datepicker"]').datepicker({
//    format: 'dd.mm.yyyy',
//    weekStart: 1,
//    startDate: date,
//    autoPick: true,
//    autoShow: true,
//    //inline: true,
//    //container: '.datepicker-inline',
//    months: ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'],
//    monthsShort: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'],
//    days: ['воскресенье', 'понедельник', 'вторник', 'среда', 'четверг', 'пятница', 'суббота'],
//    daysShort: ['вск', 'пнд', 'втр', 'срд', 'чтв', 'птн', 'сбт'],
//    daysMin: ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'],
//});

$(document).ready(function() {
    $('#buutonId').on('click', function() {
        $('#modalId').modal('open');
    });
});

$(document).ready(function () {




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

//('#my-element').datepicker([options])
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
    setInterval(function () {
        // Just move your date creation inside the interval function
        var today = new Date();
        var date = setZero(today.getDate()) + '.' + (setZero(today.getMonth() + 1)) + '.' + today.getFullYear();
        var time = setZero(today.getHours()) + ":" + setZero(today.getMinutes());// + ":" + today.getSeconds();
        var dateTime = date + " " + time; // Add the time to the date string

        document.getElementById('clock').innerHTML = dateTime;
    }, 1000);


////проверка на ввод в поле "Количество человек"
//$('.table').checkboxTable();
//$('[class^=num]').keypress(validateNumber);


//$('#datepicker').datepicker({
//    format: "dd.mm.yyyy",
//    language: "ru",
//    weekStart: 1,
//    daysOfWeekHighlighted: "0,6",
//    todayHighlight: true,
//    //autoclose: true,
//});
//$('#datepicker').datepicker("setDate", new Date());

})