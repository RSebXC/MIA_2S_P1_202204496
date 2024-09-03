import React, { useState, useRef } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import '@fortawesome/fontawesome-free/css/all.min.css';

function App() {
  const [inputText, setInputText] = useState('');
  const [outputText, setOutputText] = useState('');
  const fileInputRef = useRef(null);

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setInputText(e.target.result);
      };
      reader.readAsText(file);
    }
  };

  const ejecutarComandos = () => {
    fetch('http://localhost:8080/api/message', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ inputText }),
    })
    .then(response => response.json())
    .then(data => {
      setOutputText(data.message);
    })
    .catch(error => {
      console.error('Error:', error);
      setOutputText('Error al enviar la solicitud');
    });
  };

  return (
    <div className="App" style={{ fontFamily: 'Arial, sans-serif', backgroundColor: '#f8f9fa', height: '100vh', padding: '20px' }}>
      <header className="App-header">
        <div className="container">
          <h2 className="mb-4">Frontend</h2>
          <div className="form-group">
            <label htmlFor="inputArea">Entrada:</label>
            <textarea
              id="inputArea"
              className="form-control"
              placeholder="Ingrese los comandos aquí..."
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              rows="5"
              style={{ resize: 'none', marginBottom: '20px' }}
            />
          </div>
          <input
            type="file"
            id="fileInput"
            style={{ display: 'none' }}
            ref={fileInputRef}
            onChange={handleFileChange}
          />
          <div className="d-flex align-items-center mb-3">
            <button
              className="btn btn-primary mr-2"
              onClick={() => fileInputRef.current.click()}
              style={{ marginRight: '10px' }}
            >
              <i className="fas fa-upload"></i> Cargar Archivo
            </button>
            <button className="btn btn-success" onClick={ejecutarComandos}>
              <i className="fas fa-play"></i> Ejecutar
            </button>
          </div>
          <div className="form-group">
            <label htmlFor="outputArea">Salida:</label>
            <div
              id="outputArea"
              className="p-3"
              style={{ backgroundColor: '#000', color: '#0f0', minHeight: '100px', borderRadius: '4px' }}
            >
              {outputText ? outputText : <div># Aca se veran todos los mensajes de la ejecucion</div>}
            </div>
          </div>
        </div>
      </header>
    </div>
  );
}

export default App;
