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
                    type="button"
                    aria-label="Notifications"
                >
                    🔔
                </button>


                <button
                    type="button"
                    aria-label="Profile"
                >
                    👤
                </button>

            </div>

        </header>
    );
}
