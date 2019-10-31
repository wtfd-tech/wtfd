export function registerBackdropClickHandler(dlg: HTMLDialogElement) {
  Array.prototype.slice.call(document.getElementsByClassName("backdrop"))
    .forEach(function(elem: Element) {
      elem.addEventListener("click", function() {
        history.pushState({ foo: "baar" }, "index", "/");
        dlg.close();
      });
    });
}
export function showDialog(dlg: HTMLDialogElement) {
  dlg.showModal();
  registerBackdropClickHandler(dlg);
}
