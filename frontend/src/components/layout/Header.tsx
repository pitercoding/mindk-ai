export default function Header() {
    return (
        <header className="app-header">

            <div className="search-container">

                <input
                    type="text"
                    placeholder="Search your notes, links..."
                />

            </div>


            <div className="header-actions">

                <button
                    className="notification-button"
                    type="button"
                    aria-label="Notifications"
                >
                    🔔
                </button>


                <button
                    className="profile-button"
                    type="button"
                    aria-label="Profile"
                >
                    👤
                </button>

            </div>

        </header>
    );
}
