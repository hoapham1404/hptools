import { useStatus, useProcesses, useWindowControl } from './hooks';
import { ProcessSelector, WindowControls, StatusDisplay } from './components';

function App() {
  const { status, setStatus } = useStatus();
  const {
    processes,
    selectedProcess,
    loading: processesLoading,
    debugMode,
    setSelectedProcess,
    setDebugMode,
    fetchProcesses,
  } = useProcesses(setStatus);

  const {
    dimensions,
    currentWindowInfo,
    loading: windowLoading,
    setDimensions,
    setWindowSize,
    setWindowPosition,
    getWindowInfo,
    clearWindowInfo,
  } = useWindowControl(setStatus);

  const handleSetWindowSize = () => {
    if (selectedProcess) {
      setWindowSize(selectedProcess);
    }
  };

  const handleSetWindowPosition = () => {
    if (selectedProcess) {
      setWindowPosition(selectedProcess);
    }
  };

  const handleGetWindowInfo = () => {
    if (selectedProcess) {
      getWindowInfo(selectedProcess);
    }
  };

  const handleProcessSelect = (process: typeof selectedProcess) => {
    setSelectedProcess(process);
    // Clear current window info when selecting a new process
    if (process !== selectedProcess) {
      clearWindowInfo();
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Windows Control Panel</h1>
        
        <ProcessSelector
          processes={processes}
          selectedProcess={selectedProcess}
          loading={processesLoading}
          debugMode={debugMode}
          onProcessSelect={handleProcessSelect}
          onRefresh={fetchProcesses}
          onDebugModeChange={setDebugMode}
        />

        {selectedProcess && (
          <WindowControls
            selectedProcess={selectedProcess}
            dimensions={dimensions}
            currentWindowInfo={currentWindowInfo}
            loading={windowLoading}
            onDimensionsChange={setDimensions}
            onSetWindowSize={handleSetWindowSize}
            onSetWindowPosition={handleSetWindowPosition}
            onGetWindowInfo={handleGetWindowInfo}
          />
        )}

        <StatusDisplay status={status} />
      </div>
    </div>
  );
}

export default App;
