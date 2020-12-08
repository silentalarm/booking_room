//Highlighting rows

(function ($) {
    $.fn.checkboxTable = function () {
        target = this;

        //click on checkbox.
        $(target).on('click', 'tbody :checkbox', function () {
            var parents = $(this).parents('table');
            if ($(this).is(':checked')) {
                $(this).parents('tr').addClass('bg-select');
                //$(this).parents('tr').find('input').addClass('bg-green text-light');
            } else {
                $(this).parents('tr').removeClass('bg-select');
                $(this).parents('tr').find('input').removeClass('bg-green text-light');
                if ($(parents).find('tbody :checkbox:checked').length == 0) {
                    $(parents).find('thead :checkbox').prop('checked', false);
                }
            }
        });
    };
})(jQuery);

//защита от ввода букв в поле для количества человек
function validateNumber(event) {
    var ASCIICode = (event.which) ? event.which : event.keyCode
    if (ASCIICode > 31 && (ASCIICode < 48 || ASCIICode > 57))
        return false;
    return true;
}




$(document).ready(function(){
    //проверка на ввод в поле "Количество человек"
    $('[class^=num]').keypress(validateNumber);

    //liveTime
    setInterval(function() {
        // Just move your date creation inside the interval function
        var today = new Date();
        var date = today.getFullYear() + '.' + (today.getMonth() + 1) + '.' + today.getDate();
        var time = today.getHours() + ":" + today.getMinutes() + ":" + today.getSeconds();
        var dateTime = date + " " + time; // Add the time to the date string

        document.getElementById('clock').innerHTML = dateTime;
    }, 1000);



////проверка на checkbox
  //  const checkbox = document.querySelector('.check_box__status');
  //  const button = document.querySelector('#btn');
//
  //  checkbox.addEventListener('change', () => {
  //      if ( checkbox.checked ) {
  //          button.removeAttribute('disabled');
  //      } else {
  //          button.setAttribute('disabled', 'true');
  //      }
  //  });
//

})