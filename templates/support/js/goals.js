document.addEventListener('DOMContentLoaded', function () {
    const banana_image = document.querySelectorAll('.banana_image');
    const newAmount = document.getElementById('newAmount');
    // var targetAmount = {{.PiggyBank.TargetAmount}}
    
    // func that changes opacity
    function change_opacity() {

        var opacity_new = parseFloat(newAmount.value);
        var opacity_decimal = (opacity_new / targetAmount);
  
        banana_image[0].style.opacity = 1 - opacity_decimal;
        banana_image[1].style.opacity = opacity_decimal;
    }
  
    // change_opacity(); 
    // newAmount.addEventListener('change', change_opacity);
});