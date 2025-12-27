
// Supabase Authentication Logic

let supabaseClient;

function initSupabase() {
    if (window.AUTH_CONFIG && window.AUTH_CONFIG.enabled) {
        if (window.supabase) {
            supabaseClient = window.supabase.createClient(window.AUTH_CONFIG.supabaseUrl, window.AUTH_CONFIG.supabaseKey);
            console.log("Supabase initialized");
        } else {
            console.error("Supabase SDK not loaded");
        }
    }
}

async function handleLogin(providerId) {
    if (!supabaseClient) {
        alert("Authentication is not enabled or Supabase is not initialized.");
        return;
    }

    try {
        const { data, error } = await supabaseClient.auth.signInWithOAuth({
            provider: providerId,
            options: {
                redirectTo: window.location.origin + window.location.pathname
            }
        });

        if (error) {
            console.error("Error logging in:", error.message);
            alert("Login failed: " + error.message);
        } else {
            console.log("Login initiated:", data);
        }
    } catch (err) {
        console.error("Unexpected error:", err);
        alert("An unexpected error occurred.");
    }
}

// Initialize on load
document.addEventListener('DOMContentLoaded', initSupabase);
