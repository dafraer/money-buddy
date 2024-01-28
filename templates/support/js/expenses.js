document.addEventListener('DOMContentLoaded', function () {
    // expensess page

    // Get the modal
    var plusModal = document.getElementById("plusModal");    
    var minModal = document.getElementById("minModal");    
    // Get the button that opens the modal
    var open_plusModal = document.getElementById("open_plusModal");
    var open_minModal = document.getElementById("open_minModal");
    // Get the <span> element that closes the modal
    var close_plusModal = document.getElementById("close_plusModal");
    var close_minModal = document.getElementById("close_minModal");

    // When the user clicks the button, open the modal 
    open_plusModal.onclick = function() {
        plusModal.style.display = "block"; }
    open_minModal.onclick = function() {
        minModal.style.display = "block"; }

    // When the user clicks on <span> (x), close the modal
    close_plusModal.onclick = function() {
        plusModal.style.display = "none";}
    close_minModal.onclick = function() {
        minModal.style.display = "none";}
    // When the user clicks anywhere outside of the modal, close it
    window.onclick = function(event) {
        if (event.target == plusModal) {
        plusModal.style.display = "none";}
        if (event.target == minModal) {
            minModal.style.display = "none";}
    }
});