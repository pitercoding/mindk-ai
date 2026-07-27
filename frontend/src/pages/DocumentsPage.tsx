import DocumentList from "@/components/documents/DocumentList";


export default function DocumentsPage() {

    return (
        <section className="documents-page">

            <header className="page-header">

                <div>
                    <h1>Documents</h1>

                    <p>Manage your knowledge documents.</p>
                </div>

            </header>

            <DocumentList />

        </section>
    );
}