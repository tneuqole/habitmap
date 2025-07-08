export function editHabitMenu() {
  const btn = document.getElementById("edit-habit-btn");
  const menu = document.getElementById("edit-habit-dropdown");

  if (!btn || !menu) return;

  const closeMenu = () => {
    menu.hidden = true;
  };
  const toggleMenu = () => {
    menu.hidden = !menu.hidden;
  };

  closeMenu();

  btn.addEventListener("click", function (e) {
    e.stopPropagation();
    toggleMenu();
  });

  document.addEventListener("click", function (e) {
    if (!menu.contains(e.target) && !btn.contains(e.target)) {
      closeMenu();
    }
  });
}
