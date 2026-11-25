import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { Link, NavLink, useNavigate } from "react-router-dom";

import ConfirmDialog from "@/components/common/ConfirmDialog";
import {
    DocumentsIcon,
    HomeIcon,
    LinksIcon,
    NotesIcon,
    PdfIcon,
    SettingsIcon,
} from "@/components/layout/icons";

import { useChatSession } from "@/context/ChatSessionContext";

const menuItems = [
    {
        label: "Dashboard",
        path: "/",
        icon: <HomeIcon />,
        end: true,
    },
    {
        label: "All Notes",
        path: "/notes",
        icon: <NotesIcon />,
    },
    {
        label: "Documents",
        path: "/documents",
        icon: <DocumentsIcon />,
    },
    {
        label: "Links",
        icon: <LinksIcon />,
        soon: true as const,
    },
    {
        label: "PDFs",
        icon: <PdfIcon />,
        soon: true as const,
    },
    {
        label: "Settings",
        icon: <SettingsIcon />,
        soon: true as const,
    },
];

interface MenuState {
    sessionId: number;
    top: number;
    left: number;
}

interface SidebarProps {
    isOpen: boolean;
    onClose: () => void;
}

export default function Sidebar({ isOpen, onClose }: SidebarProps) {

    const navigate = useNavigate();

    const {
        sessions,
        isLoadingSessions,
        sessionsError,
        currentSessionId,
        newChat,
        selectSession,
        renameSession,
        deleteSession,
    } = useChatSession();

    const [menu, setMenu] = useState<MenuState | null>(null);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [editingTitle, setEditingTitle] = useState("");
    const [actionError, setActionError] = useState<string | null>(null);
    const [pendingDeleteId, setPendingDeleteId] = useState<number | null>(null);
    const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null);

    useEffect(() => {

        if (!menu) {
            return;
        }

        function closeMenu() {
            setMenu(null);
        }

        document.addEventListener("mousedown", closeMenu);
        window.addEventListener("resize", closeMenu);
        window.addEventListener("scroll", closeMenu, true);

        return () => {
            document.removeEventListener("mousedown", closeMenu);
            window.removeEventListener("resize", closeMenu);
            window.removeEventListener("scroll", closeMenu, true);
        };

    }, [menu]);

    function handleNewChat() {
        setMenu(null);
        newChat();
        navigate("/");
        onClose();
    }

    async function handleSelect(id: number) {
        setMenu(null);
        await selectSession(id);
        navigate("/");
        onClose();
    }

    function toggleMenu(
        event: React.MouseEvent<HTMLButtonElement>,
        sessionId: number,
    ) {

        if (menu?.sessionId === sessionId) {
            setMenu(null);
            return;
        }

        const rect = event.currentTarget.getBoundingClientRect();
        const menuWidth = 120;

        setMenu({
            sessionId,
            top: rect.bottom + 4,
            left: Math.min(
                rect.right - menuWidth,
                window.innerWidth - menuWidth - 8,
            ),
        });
    }

    function startRename(id: number, title: string) {
        setMenu(null);
        setEditingId(id);
        setEditingTitle(title);
    }

    async function commitRename(id: number) {

        const title = editingTitle.trim();

        setEditingId(null);

        if (!title) {
            return;
        }

        setActionError(null);

        try {
            await renameSession(id, title);
        } catch (error) {
            console.error("Failed to rename conversation:", error);
            setActionError("Unable to rename conversation.");
        }
    }

    function requestDelete(id: number) {
        setMenu(null);
        setConfirmDeleteId(id);
    }

    async function handleDelete(id: number) {

        setConfirmDeleteId(null);
        setActionError(null);
        setPendingDeleteId(id);

        try {
            await deleteSession(id);
        } catch (error) {
            console.error("Failed to delete conversation:", error);
            setActionError("Unable to delete conversation.");
        } finally {
            setPendingDeleteId(null);
        }
    }

    const confirmSession =
        confirmDeleteId !== null
            ? sessions.find(session => session.id === confirmDeleteId)
            : undefined;

    const activeMenuSession =
        menu && sessions.find(session => session.id === menu.sessionId);

    return (
        <aside
            id="app-sidebar"
            className={isOpen ? "sidebar sidebar-open" : "sidebar"}
        >

            <header className="sidebar-header">
                <h1>
                    <Link to="/" className="sidebar-logo" aria-label="MindK AI - go to dashboard">
                        🧠 MindK AI
                    </Link>
                </h1>

                <button
                    type="button"
                    className="sidebar-close-button"
                    aria-label="Close menu"
                    onClick={onClose}
                >
                    ✕
                </button>
            </header>


            <nav className="sidebar-menu">

                {menuItems.map((item) => (

                    item.soon ? (

                        <div
                            key={item.label}
                            className="sidebar-link sidebar-link-disabled"
                            aria-disabled="true"
                            title="Coming soon"
                        >
                            {item.icon}

                            <span>
                                {item.label}
                            </span>

                            <small>
                                Soon
                            </small>
                        </div>

                    ) : (

                        <NavLink
                            key={item.label}
                            to={item.path}
                            end={item.end}
                            className={({ isActive }) =>
                                isActive
                                    ? "sidebar-link active"
                                    : "sidebar-link"
                            }
                        >
                            {item.icon}

                            <span>
                                {item.label}
                            </span>
                        </NavLink>

                    )
                ))}

            </nav>


            <div className="sidebar-conversations">

                <div className="sidebar-conversations-header">
                    <span className="sidebar-section-title">
                        Conversations
                    </span>
                </div>

                <button
                    className="sidebar-new-chat-button"
                    onClick={handleNewChat}
                >
                    + New Chat
                </button>

                {actionError && (
                    <p className="sidebar-status sidebar-status-error">
                        {actionError}
                    </p>
                )}

                {isLoadingSessions ? (

                    <p className="sidebar-status">
                        Loading conversations...
                    </p>

                ) : sessionsError ? (

                    <p className="sidebar-status sidebar-status-error">
                        {sessionsError}
                    </p>

                ) : sessions.length === 0 ? (

                    <p className="sidebar-status">
                        No conversations yet.
                    </p>

                ) : (

                <ul className="conversation-list">

                    {sessions.map((session) => {

                        const isDeleting = pendingDeleteId === session.id;

                        return (
                        <li
                            key={session.id}
                            className={[
                                "conversation-item",
                                session.id === currentSessionId ? "active" : "",
                                isDeleting ? "conversation-item-deleting" : "",
                            ].filter(Boolean).join(" ")}
                        >
                            {isDeleting ? (

                                <span className="conversation-title">
                                    Deleting...
                                </span>

                            ) : editingId === session.id ? (

                                <input
                                    autoFocus
                                    className="conversation-rename-input"
                                    value={editingTitle}
                                    onChange={(event) =>
                                        setEditingTitle(event.target.value)
                                    }
                                    onBlur={() => commitRename(session.id)}
                                    onKeyDown={(event) => {

                                        if (event.key === "Enter") {
                                            event.currentTarget.blur();
                                        }

                                        if (event.key === "Escape") {
                                            setEditingId(null);
                                        }
                                    }}
                                />

                            ) : (

                                <button
                                    className="conversation-title"
                                    onClick={() => handleSelect(session.id)}
                                >
                                    {session.title}
                                </button>

                            )}

                            <button
                                className="conversation-menu-toggle"
                                disabled={isDeleting}
                                onMouseDown={(event) => event.stopPropagation()}
                                onClick={(event) => toggleMenu(event, session.id)}
                            >
                                ⋮
                            </button>
                        </li>
                        );
                    })}

                </ul>

                )}

            </div>


            <footer className="sidebar-footer">

                <button onClick={() => { navigate("/documents"); onClose(); }}>
                    ⬆ Upload
                </button>

            </footer>

            {confirmSession && (
                <ConfirmDialog
                    title="Delete conversation?"
                    message={`"${confirmSession.title}" and all its messages will be permanently deleted.`}
                    onConfirm={() => handleDelete(confirmSession.id)}
                    onCancel={() => setConfirmDeleteId(null)}
                />
            )}

            {menu && activeMenuSession && createPortal(

                <div
                    className="conversation-menu"
                    style={{ top: menu.top, left: menu.left }}
                    onMouseDown={(event) => event.stopPropagation()}
                >

                    <button
                        onClick={() =>
                            startRename(activeMenuSession.id, activeMenuSession.title)
                        }
                    >
                        Rename
                    </button>

                    <button
                        onClick={() => requestDelete(activeMenuSession.id)}
                    >
                        Delete
                    </button>

                </div>,

                document.body,
            )}

        </aside>
    );
}
