import {
    createContext,
    useContext,
    useState,
    type ReactNode,
} from "react";

import type { Note } from "@/types/note";

interface SelectedNoteContextType {
    selectedNote: Note | null;
    setSelectedNote: (note: Note | null) => void;
}

const SelectedNoteContext = createContext<
    SelectedNoteContextType | undefined
>(undefined);

interface SelectedNoteProviderProps {
    children: ReactNode;
}

export function SelectedNoteProvider({
    children,
}: SelectedNoteProviderProps) {

    const [selectedNote, setSelectedNote] =
        useState<Note | null>(null);

    return (
        <SelectedNoteContext.Provider
            value={{
                selectedNote,
                setSelectedNote,
            }}
        >
            {children}
        </SelectedNoteContext.Provider>
    );
}

export function useSelectedNote() {

    const context = useContext(
        SelectedNoteContext,
    );

    if (!context) {
        throw new Error(
            "useSelectedNote must be used inside SelectedNoteProvider"
        );
    }

    return context;
}
