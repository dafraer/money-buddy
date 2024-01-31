document.addEventListener('DOMContentLoaded', function () {
    // load header and footer in every file
    $(function(){
        $("#includedheader").load("support/header.html");
        $("#includedheaderacc").load("support/headeracc.html");
        $("#includedfooter").load("support/footer.html"); 
    });

    // function that makes text visible by click
    function toggleVisibility() {
        var hiddenText = document.getElementById('hiddenText');
        hiddenText.style.display = (hiddenText.style.display === 'none' || hiddenText.style.display === '') ? 'inline' : 'none';}
 
});