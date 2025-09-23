import { useState, useEffect } from 'react';
import { WindowService, ProcessInfo, WindowInfo } from '../bindings/hptools';

function App() {
  const [processes, setProcesses] = useState<ProcessInfo[]>([]);
  const [selectedProcess, setSelectedProcess] = useState<ProcessInfo | null>(null);
  const [width, setWidth] = useState<number>(800);
  const [height, setHeight] = useState<number>(600);
  const [x, setX] = useState<number>(100);
  const [y, setY] = useState<number>(100);
  const [currentWindowInfo, setCurrentWindowInfo] = useState<WindowInfo | null>(null);
  const [status, setStatus] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [debugMode, setDebugMode] = useState<boolean>(false);

  // Fetch application processes on component mount
  useEffect(() => {
    fetchProcesses();
  }, [debugMode]); // Re-fetch when debug mode changes

  const fetchProcesses = async () => {
    try {
      setLoading(true);
      const apps = debugMode 
        ? await WindowService.GetAllProcessesWithWindows()
        : await WindowService.GetApplicationProcesses();
      setProcesses(apps);
      setStatus(`Found ${apps.length} ${debugMode ? 'processes with windows' : 'application processes'}`);
    } catch (error) {
      console.error('Error fetching processes:', error);
      setStatus(`Error: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const handleSetWindowSize = async () => {
    if (!selectedProcess) {
      setStatus('Please select a process first');
      return;
    }

    try {
      setLoading(true);
      await WindowService.SetWindowSize(selectedProcess.pid, width, height);
      setStatus(`✅ Set window size to ${width}x${height} for ${selectedProcess.imageName}`);
      // Refresh window info after resize
      getWindowInfo();
    } catch (error) {
      console.error('Error setting window size:', error);
      setStatus(`❌ Error: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const handleSetWindowPosition = async () => {
    if (!selectedProcess) {
      setStatus('Please select a process first');
      return;
    }

    try {
      setLoading(true);
      await WindowService.SetWindowPosition(selectedProcess.pid, x, y, width, height);
      setStatus(`✅ Set window position to (${x}, ${y}) and size to ${width}x${height} for ${selectedProcess.imageName}`);
      // Refresh window info after move/resize
      getWindowInfo();
    } catch (error) {
      console.error('Error setting window position:', error);
      setStatus(`❌ Error: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const getWindowInfo = async () => {
    if (!selectedProcess) {
      setStatus('Please select a process first');
      return;
    }

    try {
      setLoading(true);
      const info = await WindowService.GetWindowInfo(selectedProcess.pid);
      if (info) {
        setCurrentWindowInfo(info);
        setStatus(`📏 Current window: ${info.width}x${info.height} at position (${info.x}, ${info.y})`);
        
        // Update form fields with current values
        setWidth(info.width);
        setHeight(info.height);
        setX(info.x);
        setY(info.y);
      }
    } catch (error) {
      console.error('Error getting window info:', error);
      setStatus(`❌ Error: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Windows Control Panel</h1>
        
        {/* Process Selection */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">Select Application</h2>
          
          <div className="flex gap-4 mb-4">
            <select
              value={selectedProcess?.pid || ''}
              onChange={(e) => {
                const pid = parseInt(e.target.value);
                const process = processes.find(p => p.pid === pid);
                setSelectedProcess(process || null);
                setCurrentWindowInfo(null);
                setStatus(process ? `Selected: ${process.imageName} - "${process.windowTitle}"` : 'No process selected');
              }}
              className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select an application...</option>
              {processes.map((process) => (
                <option key={process.pid} value={process.pid}>
                  {process.imageName} - "{process.windowTitle}" (PID: {process.pid}) [{process.windowCount} windows]
                </option>
              ))}
            </select>
            
            <button
              onClick={fetchProcesses}
              disabled={loading}
              className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:opacity-50"
            >
              {loading ? '🔄' : '🔄 Refresh'}
            </button>
          </div>

          {/* Debug Mode Toggle */}
          <div className="flex items-center gap-4 mb-4">
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={debugMode}
                onChange={(e) => setDebugMode(e.target.checked)}
                className="rounded"
              />
              <span className="text-sm text-gray-600">
                Debug Mode (Show all processes with windows)
              </span>
            </label>
            
            {debugMode && (
              <div className="text-xs text-yellow-600 bg-yellow-50 px-2 py-1 rounded">
                ⚠️ Debug mode shows all processes, including system ones
              </div>
            )}
          </div>
        </div>

        {/* Window Controls */}
        {selectedProcess && (
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4">
              Control Window: {selectedProcess.imageName}
            </h2>

            {/* Current Window Info */}
            <div className="mb-6">
              <button
                onClick={getWindowInfo}
                disabled={loading}
                className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 disabled:opacity-50 mb-3"
              >
                📏 Get Current Window Info
              </button>
              
              {currentWindowInfo && (
                <div className="bg-gray-100 p-3 rounded-md text-sm">
                  <strong>Current Window:</strong> {currentWindowInfo.width}x{currentWindowInfo.height} 
                  at position ({currentWindowInfo.x}, {currentWindowInfo.y})
                </div>
              )}
            </div>

            {/* Size Controls */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Width
                </label>
                <input
                  type="number"
                  value={width}
                  onChange={(e) => setWidth(parseInt(e.target.value) || 800)}
                  min="100"
                  max="3840"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Height
                </label>
                <input
                  type="number"
                  value={height}
                  onChange={(e) => setHeight(parseInt(e.target.value) || 600)}
                  min="100"
                  max="2160"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  X Position
                </label>
                <input
                  type="number"
                  value={x}
                  onChange={(e) => setX(parseInt(e.target.value) || 100)}
                  min="0"
                  max="5000"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Y Position
                </label>
                <input
                  type="number"
                  value={y}
                  onChange={(e) => setY(parseInt(e.target.value) || 100)}
                  min="0"
                  max="5000"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex gap-4">
              <button
                onClick={handleSetWindowSize}
                disabled={loading}
                className="px-6 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:opacity-50"
              >
                📏 Resize Window
              </button>
              
              <button
                onClick={handleSetWindowPosition}
                disabled={loading}
                className="px-6 py-2 bg-purple-500 text-white rounded-md hover:bg-purple-600 disabled:opacity-50"
              >
                🎯 Move & Resize
              </button>
            </div>

            {/* Quick Size Presets */}
            <div className="mt-6">
              <h3 className="text-sm font-medium text-gray-700 mb-2">Quick Presets:</h3>
              <div className="flex gap-2 flex-wrap">
                {[
                  { name: "HD", w: 1280, h: 720 },
                  { name: "FHD", w: 1920, h: 1080 },
                  { name: "Square", w: 800, h: 800 },
                  { name: "Small", w: 640, h: 480 },
                  { name: "Wide", w: 1200, h: 600 },
                  { name: "Custom 1", w: 1860, h: 1000, x: 30, y: 30 },
                ].map((preset) => (
                  <button
                    key={preset.name}
                    onClick={() => {
                      setWidth(preset.w);
                      setHeight(preset.h);
                      if (preset.x !== undefined) setX(preset.x);
                      if (preset.y !== undefined) setY(preset.y);     
                    }}
                    className="px-3 py-1 text-xs bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
                  >
                    {preset.name} ({preset.w}x{preset.h}{preset.x !== undefined ? ` @${preset.x},${preset.y}` : ''})
                  </button>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Status Display */}
        {status && (
          <div className="bg-white rounded-lg shadow-md p-4">
            <h3 className="text-sm font-medium text-gray-700 mb-2">Status:</h3>
            <div className="text-sm text-gray-600 font-mono bg-gray-50 p-3 rounded">
              {status}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
