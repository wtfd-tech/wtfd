(function() {
  var loginButton = document.getElementById("loginbutton");
  if (loginButton != null) {
    var loginDialogCancelButton = document.getElementById("logincancelbutton");
    var loginDialog = document.getElementById("logindialog");
    dialogPolyfill.registerDialog(loginDialog);

    loginButton.addEventListener("click", function() {
      loginDialog.showModal();
    });
    loginDialogCancelButton.addEventListener("click", function() {
      loginDialog.close();
    });

    var registerButton = document.getElementById("registerbutton");
    var registerDialogCancelButton = document.getElementById(
      "registercancelbutton"
    );
    var registerDialog = document.getElementById("registerdialog");
    dialogPolyfill.registerDialog(registerDialog);

    registerButton.addEventListener("click", function() {
      loginDialog.close();
      registerDialog.showModal();
    });
    registerDialogCancelButton.addEventListener("click", function() {
      registerDialog.close();
    });
  } else {
    var logoutButton = document.getElementById("logoutbutton");
    logoutButton.addEventListener("click", function() {
      location.href = "/logout";
    });
  }
})();
