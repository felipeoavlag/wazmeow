import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { Session, SessionFilters, DashboardData } from '@/lib/types/api';

interface AppState {
  // Sessões
  sessions: Session[];
  selectedSession: Session | null;
  sessionFilters: SessionFilters;
  
  // UI State
  theme: 'light' | 'dark' | 'system';
  sidebarOpen: boolean;
  isMobile: boolean;
  loading: boolean;
  
  // Configurações
  apiUrl: string;
  refreshInterval: number;
  
  // Dashboard
  dashboardData: DashboardData | null;
  
  // Actions - Sessões
  setSessions: (sessions: Session[]) => void;
  addSession: (session: Session) => void;
  updateSession: (id: string, updates: Partial<Session>) => void;
  removeSession: (id: string) => void;
  setSelectedSession: (session: Session | null) => void;
  setSessionFilters: (filters: SessionFilters) => void;
  
  // Actions - UI
  setTheme: (theme: 'light' | 'dark' | 'system') => void;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  setIsMobile: (mobile: boolean) => void;
  setLoading: (loading: boolean) => void;
  
  // Actions - Configurações
  setApiUrl: (url: string) => void;
  setRefreshInterval: (interval: number) => void;
  
  // Actions - Dashboard
  setDashboardData: (data: DashboardData) => void;
  
  // Utility actions
  reset: () => void;
}

const initialState = {
  sessions: [],
  selectedSession: null,
  sessionFilters: {},
  theme: 'system' as const,
  sidebarOpen: true,
  isMobile: false,
  loading: false,
  apiUrl: 'http://localhost:8080',
  refreshInterval: 30000, // 30 segundos
  dashboardData: null,
};

export const useAppStore = create<AppState>()(
  persist(
    (set, get) => ({
      ...initialState,
      
      // Actions - Sessões
      setSessions: (sessions) => set({ sessions }),
      
      addSession: (session) =>
        set((state) => ({
          sessions: [...state.sessions, session],
        })),
      
      updateSession: (id, updates) =>
        set((state) => ({
          sessions: state.sessions.map((session) =>
            session.id === id ? { ...session, ...updates } : session
          ),
          selectedSession:
            state.selectedSession?.id === id
              ? { ...state.selectedSession, ...updates }
              : state.selectedSession,
        })),
      
      removeSession: (id) =>
        set((state) => ({
          sessions: state.sessions.filter((session) => session.id !== id),
          selectedSession:
            state.selectedSession?.id === id ? null : state.selectedSession,
        })),
      
      setSelectedSession: (session) => set({ selectedSession: session }),
      
      setSessionFilters: (filters) => set({ sessionFilters: filters }),
      
      // Actions - UI
      setTheme: (theme) => set({ theme }),
      
      toggleSidebar: () =>
        set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      
      setSidebarOpen: (open) => set({ sidebarOpen: open }),
      
      setIsMobile: (mobile) => set({ isMobile: mobile }),
      
      setLoading: (loading) => set({ loading }),
      
      // Actions - Configurações
      setApiUrl: (url) => set({ apiUrl: url }),
      
      setRefreshInterval: (interval) => set({ refreshInterval: interval }),
      
      // Actions - Dashboard
      setDashboardData: (data) => set({ dashboardData: data }),
      
      // Utility actions
      reset: () => set(initialState),
    }),
    {
      name: 'wazmeow-app-store',
      partialize: (state) => ({
        theme: state.theme,
        sidebarOpen: state.sidebarOpen,
        apiUrl: state.apiUrl,
        refreshInterval: state.refreshInterval,
        sessionFilters: state.sessionFilters,
      }),
    }
  )
);

// Selectors para otimização
export const useSessionsStore = () => {
  const sessions = useAppStore((state) => state.sessions);
  const selectedSession = useAppStore((state) => state.selectedSession);
  const sessionFilters = useAppStore((state) => state.sessionFilters);
  const setSessions = useAppStore((state) => state.setSessions);
  const addSession = useAppStore((state) => state.addSession);
  const updateSession = useAppStore((state) => state.updateSession);
  const removeSession = useAppStore((state) => state.removeSession);
  const setSelectedSession = useAppStore((state) => state.setSelectedSession);
  const setSessionFilters = useAppStore((state) => state.setSessionFilters);
  
  return {
    sessions,
    selectedSession,
    sessionFilters,
    setSessions,
    addSession,
    updateSession,
    removeSession,
    setSelectedSession,
    setSessionFilters,
  };
};

export const useUIStore = () => {
  const theme = useAppStore((state) => state.theme);
  const sidebarOpen = useAppStore((state) => state.sidebarOpen);
  const isMobile = useAppStore((state) => state.isMobile);
  const loading = useAppStore((state) => state.loading);
  const setTheme = useAppStore((state) => state.setTheme);
  const toggleSidebar = useAppStore((state) => state.toggleSidebar);
  const setSidebarOpen = useAppStore((state) => state.setSidebarOpen);
  const setIsMobile = useAppStore((state) => state.setIsMobile);
  const setLoading = useAppStore((state) => state.setLoading);
  
  return {
    theme,
    sidebarOpen,
    isMobile,
    loading,
    setTheme,
    toggleSidebar,
    setSidebarOpen,
    setIsMobile,
    setLoading,
  };
};

export const useConfigStore = () => {
  const apiUrl = useAppStore((state) => state.apiUrl);
  const refreshInterval = useAppStore((state) => state.refreshInterval);
  const setApiUrl = useAppStore((state) => state.setApiUrl);
  const setRefreshInterval = useAppStore((state) => state.setRefreshInterval);
  
  return {
    apiUrl,
    refreshInterval,
    setApiUrl,
    setRefreshInterval,
  };
};