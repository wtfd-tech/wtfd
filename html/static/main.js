(function() {
  var dialog = document.querySelector("dialog");
  dialogPolyfill.registerDialog(dialog);
  var loginButton = document.getElementById("loginbutton");
  var loginDialog = document.getElementById("logindialog");
  loginButton.addEventListener("click", function() {
    loginDialog.showModal();
  });
})();
