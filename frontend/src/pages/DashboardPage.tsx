import { useEffect, useState } from "react";

import ChatPanel from "@/components/dashboard/ChatPanel";
import CurrentNote from "@/components/dashboard/CurrentNote";
import KnowledgeBase from "@/components/dashboard/KnowledgeBase";

import { useSelectedNote } from "@/context/SelectedNoteContext";

import { getNotes } from "@/services/noteService";

import type { Note } from "@/types/note";

export default function DashboardPage() {

    const [notes, setNotes] = useState<Note[]>([]);

    const {
        selectedNote,
        setSelectedNote,
    } = useSelectedNote();


    async function loadNotes() {

        try {

            const response = await getNotes();

            setNotes(response);

            if (response.length > 0) {
                setSelectedNote(response[0]);
            }

        } catch (error) {

            console.error(
                "Failed to load notes:",
                error,
            );
        }
    }


    useEffect(() => {
        loadNotes();
    }, []);


    return (
        <section className="dashboard-page">

            <div className="dashboard-grid">

                <KnowledgeBase
                    notes={notes}
                />

                <CurrentNote
                    note={selectedNote}
                />

                <ChatPanel />

            </div>

        </section>
    );
}
