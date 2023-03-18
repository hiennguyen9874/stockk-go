import TVChartContainer from 'components/common/TVChartContainer';
import Datafeed from 'containers/DataFeed';
import { API_DATAFEED_URL, API_STORAGE_URL } from 'configs/api-server';

import './App.css';

function App() {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-explicit-any
  return (
    <div className="App">
      <div className="App-body">
        <TVChartContainer
          datafeed={Datafeed}
          symbol="TCB"
          datafeedUrl={API_DATAFEED_URL}
          // chartsStorageUrl={API_STORAGE_URL}
        />
      </div>
    </div>
  );
}

export default App;
