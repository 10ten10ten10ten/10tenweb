$(document).ready(function () {
    // Theme Toggling Logic
    const themeToggleBtn = $('#theme-toggle');
    const iconSun = themeToggleBtn.find('.fa-sun');
    const iconMoon = themeToggleBtn.find('.fa-moon');

    // Check for saved theme preference or use system preference
    const currentTheme = localStorage.getItem('theme') ||
        (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');

    // Apply initial theme
    if (currentTheme === 'dark') {
        $('html').attr('data-theme', 'dark');
        iconSun.show();
        iconMoon.hide();
    } else {
        $('html').attr('data-theme', 'light');
        iconSun.hide();
        iconMoon.show();
    }

    // Toggle event
    themeToggleBtn.on('click', function () {
        const isDark = $('html').attr('data-theme') === 'dark';

        if (isDark) {
            $('html').attr('data-theme', 'light');
            localStorage.setItem('theme', 'light');
            iconSun.hide();
            iconMoon.show();
        } else {
            $('html').attr('data-theme', 'dark');
            localStorage.setItem('theme', 'dark');
            iconSun.show();
            iconMoon.hide();
        }
    });


    // Add slight animation to body opacity on load
    $('body').css('opacity', '1');
});
