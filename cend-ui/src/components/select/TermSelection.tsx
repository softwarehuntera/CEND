import { TermEntry } from '@/types/backendTypes';

interface TermSelectionProps {
    selectedTerm: TermEntry | null;
    selectedTermCluster: TermEntry[];

}

export default function TermSelection({ selectedTerm, selectedTermCluster }: TermSelectionProps) {
    if (!selectedTerm) {
        return <p>No term selected.</p>;
    }

    return (
        <div>
            <h2>Selected Term</h2>
            <div style={{ marginBottom: '20px', padding: '10px', border: '1px solid #ccc', borderRadius: '5px' }}>
                <p><strong>Document:</strong> {selectedTerm.document}</p>
                <p><strong>ID:</strong> {selectedTerm.id}</p>
                <p><strong>Is Preferred:</strong> {selectedTerm.isPreferred ? 'Yes' : 'No'}</p>
                {selectedTerm.fields && (
                    <div>
                        <strong>Fields:</strong>
                        <ul>
                            {Object.entries(selectedTerm.fields).map(([key, value]) => (
                                <li key={key}>
                                    <strong>{key}:</strong> {value}
                                </li>
                            ))}
                        </ul>
                    </div>
                )}
                <p><strong>Preferred Documents:</strong> {selectedTerm.preferredDocuments.join(', ') || 'None'}</p>
            </div>

            <h2>Cluster</h2>
            {selectedTermCluster.length > 0 ? (
                <ul>
                    {selectedTermCluster.map((term) => (
                        <li
                            key={term.id}
                            style={{
                                marginBottom: '10px',
                                padding: '10px',
                                border: '1px solid #ccc',
                                borderRadius: '5px',
                            }}
                        >
                            <p><strong>Document:</strong> {term.document}</p>
                            <p><strong>ID:</strong> {term.id}</p>
                            <p><strong>Is Preferred:</strong> {term.isPreferred ? 'Yes' : 'No'}</p>
                        </li>
                    ))}
                </ul>
            ) : (
                <p>No terms in the cluster.</p>
            )}
        </div>
    );
}