import type { Note } from "@/types/note";

interface NoteListProps {
    notes: Note[];
}

export default function NoteList({
    notes,
}: NoteListProps) {

    return (
        <section>
            {notes.map((note) => (
                <article key={note.id}>
                    <h3>
                        {note.title}
                    </h3>

                    <p>
                        {note.content}
                    </p>
                </article>
            ))}
        </section>
    );
}