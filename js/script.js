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


    // Load Configuration
    const data = window.CONFIG || {};
    const config = data.company || {};
    const siteTitle = data.title;

    // Populate elements with data-config attribute (checking both config and top-level data)
    $('[data-config]').each(function () {
        const key = $(this).data('config');
        if (config[key]) {
            $(this).text(config[key]);
        } else if (data[key]) {
            $(this).text(data[key]);
        }
    });

    // Set Document Title
    if (siteTitle) {
        document.title = siteTitle;
    } else if (config.companyShortName) {
        document.title = config.companyShortName + " - Static";
    }

    // Set Year from config if available
    if (config.companyYear) {
        $('#year').text(config.companyYear);
    }

    // Add slight animation to body opacity on load
    $('body').css('opacity', '1');
});
