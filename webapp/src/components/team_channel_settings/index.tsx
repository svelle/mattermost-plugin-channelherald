import React, {useState, useEffect, useCallback} from 'react';

type TeamConfig = {
    channel_name: string;
    enabled: boolean;
    show_purpose: boolean;
    notify_on_conversion: boolean;
};

type Team = {
    id: string;
    display_name: string;
    name: string;
};

type Props = {
    id: string;
    value: string;
    disabled: boolean;
    onChange: (id: string, value: string) => void;
    setSaveNeeded: () => void;
};

const defaultTeamConfig = (): TeamConfig => ({
    channel_name: '',
    enabled: false,
    show_purpose: false,
    notify_on_conversion: false,
});

const parseConfigs = (value: string): Record<string, TeamConfig> => {
    if (!value) {
        return {};
    }
    try {
        return JSON.parse(value) as Record<string, TeamConfig>;
    } catch {
        return {};
    }
};

const styles: Record<string, React.CSSProperties> = {
    teamCard: {
        marginBottom: '16px',
        padding: '16px',
        border: '1px solid rgba(var(--center-channel-color-rgb, 63, 67, 80), 0.16)',
        borderRadius: '4px',
        backgroundColor: 'var(--center-channel-bg, #fff)',
    },
    teamTitle: {
        margin: '0 0 12px 0',
        fontSize: '14px',
        fontWeight: 600,
        color: 'var(--center-channel-color, #3f4350)',
    },
    row: {
        marginBottom: '10px',
        display: 'flex',
        flexDirection: 'column' as const,
        gap: '6px',
    },
    checkRow: {
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
        marginBottom: '8px',
        cursor: 'pointer',
        fontSize: '14px',
        color: 'var(--center-channel-color, #3f4350)',
    },
    helpText: {
        fontSize: '12px',
        color: 'var(--center-channel-color-56, rgba(63,67,80,0.56))',
        marginTop: '2px',
    },
    label: {
        fontSize: '13px',
        fontWeight: 500,
        color: 'var(--center-channel-color, #3f4350)',
        marginBottom: '4px',
    },
    subSettings: {
        paddingLeft: '16px',
        borderLeft: '2px solid rgba(var(--center-channel-color-rgb, 63, 67, 80), 0.08)',
        marginTop: '8px',
    },
};

const TeamChannelSettings: React.FC<Props> = ({id, value, disabled, onChange, setSaveNeeded}) => {
    const [teams, setTeams] = useState<Team[]>([]);
    const [configs, setConfigs] = useState<Record<string, TeamConfig>>(() => parseConfigs(value));
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        fetch('/api/v4/teams?per_page=200', {credentials: 'include'}).
            then((res) => {
                if (!res.ok) {
                    throw new Error(`Failed to fetch teams: ${res.status}`);
                }
                return res.json();
            }).
            then((data: Team[]) => {
                setTeams(data);
                setLoading(false);
            }).
            catch((err: Error) => {
                setError(err.message);
                setLoading(false);
            });
    }, []);

    const updateConfig = useCallback((teamId: string, field: keyof TeamConfig, val: string | boolean) => {
        setConfigs((prev) => {
            const next = {
                ...prev,
                [teamId]: {
                    ...(prev[teamId] ?? defaultTeamConfig()),
                    [field]: val,
                },
            };
            onChange(id, JSON.stringify(next));
            setSaveNeeded();
            return next;
        });
    }, [id, onChange, setSaveNeeded]);

    if (loading) {
        return <div style={{padding: '8px', color: 'var(--center-channel-color-56)'}}>{'Loading teams...'}</div>;
    }

    if (error) {
        return <div style={{padding: '8px', color: 'var(--error-text, #d24b4e)'}}>{'Error loading teams: '}{error}</div>;
    }

    if (teams.length === 0) {
        return <div style={{padding: '8px', color: 'var(--center-channel-color-56)'}}>{'No teams found.'}</div>;
    }

    return (
        <div>
            {teams.map((team) => {
                const cfg = configs[team.id] ?? defaultTeamConfig();
                return (
                    <div
                        key={team.id}
                        style={styles.teamCard}
                    >
                        <div style={styles.teamTitle}>{team.display_name}</div>

                        <label style={styles.checkRow}>
                            <input
                                type='checkbox'
                                checked={cfg.enabled}
                                disabled={disabled}
                                onChange={(e) => updateConfig(team.id, 'enabled', e.target.checked)}
                            />
                            {'Enable notifications for this team'}
                        </label>

                        {cfg.enabled && (
                            <div style={styles.subSettings}>
                                <div style={styles.row}>
                                    <div style={styles.label}>{'Notification channel name'}</div>
                                    <input
                                        type='text'
                                        className='form-control'
                                        value={cfg.channel_name}
                                        disabled={disabled}
                                        placeholder='e.g. town-square'
                                        onChange={(e) => updateConfig(team.id, 'channel_name', e.target.value)}
                                    />
                                    <div style={styles.helpText}>
                                        {'Enter the channel URL name (the slug shown in the channel URL, e.g. "town-square"). The channel must exist in this team.'}
                                    </div>
                                </div>

                                <label style={styles.checkRow}>
                                    <input
                                        type='checkbox'
                                        checked={cfg.show_purpose}
                                        disabled={disabled}
                                        onChange={(e) => updateConfig(team.id, 'show_purpose', e.target.checked)}
                                    />
                                    {'Include channel purpose in the notification'}
                                </label>

                                <label style={styles.checkRow}>
                                    <input
                                        type='checkbox'
                                        checked={cfg.notify_on_conversion}
                                        disabled={disabled}
                                        onChange={(e) => updateConfig(team.id, 'notify_on_conversion', e.target.checked)}
                                    />
                                    {'Also notify when a private channel is made public'}
                                </label>
                            </div>
                        )}
                    </div>
                );
            })}
        </div>
    );
};

export default TeamChannelSettings;
