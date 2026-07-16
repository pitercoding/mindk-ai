import { NavLink } from "react-router-dom";

const menuItems = [
    {
        label: "Dashboard",
        path: "/",
    },
    {
        label: "All Notes",
        path: "/notes",
    },
    {
        label: "Documents",
        path: "/documents",
    },
    {
        label: "Links",
        path: "/links",
    },
    {
        label: "PDFs",
        path: "/pdfs",
        badge: "Soon",
    },
    {
        label: "Settings",
        path: "/settings",
    },
];

export default function Sidebar() {
    return (
        <aside className="sidebar">

            <header className="sidebar-header">
                <h1>
                    🧠 MindK AI
                </h1>
            </header>


            <nav className="sidebar-menu">

                {menuItems.map((item) => (
                    <NavLink
                        key={item.label}
                        to={item.path}
                        className={({ isActive }) =>
                            isActive
                                ? "sidebar-link active"
                                : "sidebar-link"
                        }
                    >
                        <span>
                            {item.label}
                        </span>

                        {item.badge && (
                            <small>
                                {item.badge}
                            </small>
                        )}

                    </NavLink>
                ))}

            </nav>


            <footer className="sidebar-footer">

                <button>
                    ⬆ Upload
                </button>

            </footer>

        </aside>
    );
}
