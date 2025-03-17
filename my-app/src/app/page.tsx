'use client';

import { useState, useEffect, useRef } from 'react';

// Matrix rain characters
const matrixChars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$+-*/=%"\'#&_(),.;:?!\\|{}<>[]^~あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをん';

// Define interface for queue items
interface QueueItem {
  from: string;
  to: string;
  start: number;
  end: number;
  char?: string;
}

// Text scramble effect for glitchy text rewriting
class TextScramble {
  chars: string;
  queue: QueueItem[];
  frameRequest: number;
  frame: number;
  element: HTMLElement | null;
  resolve: Function | null;
  originalText: string;
  
  constructor() {
    this.chars = '!<>-_\\/[]{}—=+*^?#________';
    this.queue = [];
    this.frameRequest = 0;
    this.frame = 0;
    this.element = null;
    this.resolve = null;
    this.originalText = '';
  }
  
  setText(element: HTMLElement, newText: string) {
    this.element = element;
    const oldText = element.innerText;
    this.originalText = oldText;
    const length = Math.max(oldText.length, newText.length);
    
    return new Promise((resolve) => {
      this.resolve = resolve;
      
      // Reset the element state to prevent infinite growth
      element.style.fontSize = '';
      element.setAttribute('data-text', newText);
      
      this.queue = []; // Clear previous queue
      
      for (let i = 0; i < length; i++) {
        const from = oldText[i] || '';
        const to = newText[i] || '';
        const start = Math.floor(Math.random() * 20);
        const end = start + Math.floor(Math.random() * 20);
        this.queue.push({ from, to, start, end });
      }
      
      cancelAnimationFrame(this.frameRequest);
      this.frame = 0;
      this.update();
    });
  }
  
  update() {
    let output = '';
    let complete = 0;
    
    for (let i = 0, n = this.queue.length; i < n; i++) {
      let { from, to, start, end, char } = this.queue[i];
      
      if (this.frame >= end) {
        complete++;
        output += to;
      } else if (this.frame >= start) {
        if (!char || Math.random() < 0.28) {
          char = this.randomChar();
          this.queue[i].char = char;
        }
        output += `<span class="scramble-text">${char}</span>`;
      } else {
        output += from;
      }
    }
    
    if (this.element) {
      this.element.innerHTML = output;
    }
    
    if (complete === this.queue.length) {
      if (this.resolve) {
        this.resolve();
      }
    } else {
      this.frameRequest = requestAnimationFrame(() => this.update());
      this.frame++;
    }
  }
  
  randomChar() {
    return this.chars[Math.floor(Math.random() * this.chars.length)];
  }
  
  resetToOriginal(element: HTMLElement) {
    if (element && this.originalText) {
      element.innerHTML = this.originalText;
      element.setAttribute('data-text', this.originalText);
    }
  }
}

export default function Home() {
  const [email, setEmail] = useState('');
  const [submitted, setSubmitted] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);
  const [charCount, setCharCount] = useState(0);
  const [isTextGlitching, setIsTextGlitching] = useState(false);
  const [maskedEmail, setMaskedEmail] = useState('');
  const matrixRef = useRef<HTMLDivElement>(null);
  const headingRef = useRef<HTMLHeadingElement>(null);
  const emailRef = useRef<HTMLParagraphElement>(null);
  const textScrambleRef = useRef<TextScramble>(new TextScramble());
  const maxChars = 100;
  
  // Matrix rain effect with slower animation
  useEffect(() => {
    if (!matrixRef.current) return;
    
    const matrixBg = matrixRef.current;
    // Increase density of matrix columns - reduced spacing between columns
    const columns = Math.floor(window.innerWidth / 15); // Reduced from 20 to 15 for more columns
    
    // Clear any existing columns
    matrixBg.innerHTML = '';
    
    // Create columns with varied spacing
    for (let i = 0; i < columns; i++) {
      const column = document.createElement('div');
      column.className = 'matrix-column';
      
      // Add some randomness to horizontal positioning for natural look
      const horizontalOffset = Math.random() * 8 - 4; // -4px to +4px offset
      column.style.left = `${i * 15 + horizontalOffset}px`; // Reduced from 20 to 15 spacing
      
      // Make rain much slower with random duration between 7 and 15 seconds
      // This creates a more relaxed "chilling" effect with natural variation
      column.style.animationDuration = `${Math.random() * 8 + 7}s`;
      
      // Create random characters for each column with fade effect
      let content = '';
      const numChars = 60; // Increased from 50 to 60 for more characters
      // Add some randomness to character density
      const charSpacing = Math.random() < 0.3 ? 
        '<br><br>' : // 30% chance of more spacing between characters
        '<br>';      // 70% chance of normal spacing
        
      for (let j = 0; j < numChars; j++) {
        // Create variation in fade delays for more natural flow
        const fadeDelay = (j / numChars * 4) + (Math.random() * 2);
        // Add randomness to which characters are shown
        if (Math.random() > 0.1) { // 10% chance to skip a character (create gaps)
          // Increase chance of Japanese characters
          const useJapanese = Math.random() < 0.4; // 40% chance of Japanese character
          const char = useJapanese ? 
            matrixChars.slice(matrixChars.indexOf('あ')) : // Get from Japanese range
            matrixChars.slice(0, matrixChars.indexOf('あ')); // Get from Latin range
          
          const randomChar = char[Math.floor(Math.random() * char.length)];
          content += `<span style="animation-delay: ${fadeDelay}s;">${randomChar}</span>${charSpacing}`;
        } else {
          content += charSpacing;
        }
      }
      
      column.innerHTML = content;
      matrixBg.appendChild(column);
    }
  }, []);
  
  // Text rewriting effect
  useEffect(() => {
    if (!headingRef.current) return;
    
    const heading = headingRef.current;
    const textScramble = textScrambleRef.current;
    
    // Initial text
    const originalText = "ENTER YOUR EMAIL";
    // Alternate texts for scrambling
    const altTexts = [
      "CONNECT TO SYSTEM",
      "ACCESS REQUIRED",
      "IDENTITY VERIFY",
      "JOIN THE NETWORK",
      "ENTER YOUR EMAIL"
    ];
    
    let isActive = true;
    
    // Periodically scramble text
    const interval = setInterval(() => {
      if (!isActive) return;
      
      const randomText = altTexts[Math.floor(Math.random() * altTexts.length)];
      textScramble.setText(heading, randomText)
        .then(() => {
          if (!isActive) return;
          
          // After a delay, reset to original text
          setTimeout(() => {
            if (!isActive) return;
            textScramble.setText(heading, originalText);
          }, 2000);
        });
    }, 8000);
    
    return () => {
      isActive = false;
      clearInterval(interval);
      textScramble.resetToOriginal(heading);
    };
  }, []);
  
  // Add scanline effect
  useEffect(() => {
    const scanline = document.createElement('div');
    scanline.className = 'scanline';
    document.body.appendChild(scanline);
    
    return () => {
      document.body.removeChild(scanline);
    };
  }, []);

  // Glitching email effect
  useEffect(() => {
    if (!submitted || !emailRef.current) return;
    
    // Create a masked version of the email with special characters
    const specialChars = '*^%$#@!?';
    let emailGlitchInterval: NodeJS.Timeout;
    
    const updateMaskedEmail = () => {
      let result = '';
      
      // Find the position of the @ symbol
      const atSymbolIndex = email.indexOf('@');
      
      // Process the email to mask random parts
      for (let i = 0; i < email.length; i++) {
        // First two characters are always visible
        if (i < 2) {
          result += email[i];
        }
        // One character before @ is always visible
        else if (atSymbolIndex > 0 && i === atSymbolIndex - 1) {
          result += email[i];
        }
        // Middle part of username has normal masking (show with 60% probability)
        else if (i > 1 && i < atSymbolIndex - 1) {
          if (Math.random() > 0.4) {
            result += email[i];
          } else {
            const randomChar = specialChars[Math.floor(Math.random() * specialChars.length)];
            result += randomChar;
          }
        }
        // Domain part (including @) always changes
        else {
          const randomChar = specialChars[Math.floor(Math.random() * specialChars.length)];
          result += randomChar;
        }
      }
      setMaskedEmail(result);
    };
    
    // Initial update
    updateMaskedEmail();
    
    // Update the masked email periodically for glitch effect
    emailGlitchInterval = setInterval(() => {
      updateMaskedEmail();
    }, 150); // Fast updates for glitch effect
    
    return () => {
      clearInterval(emailGlitchInterval);
    };
  }, [submitted, email]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (email) {
      setShowSuccess(true);
      setTimeout(() => {
        setSubmitted(true);
        setShowSuccess(false);
      }, 1000);
    }
  };
  
  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    if (value.length <= maxChars) {
      setEmail(value);
      setCharCount(value.length);
    }
  };
  
  const restartForm = () => {
    setEmail('');
    setCharCount(0);
    setSubmitted(false);
  };

  return (
    <>
      <div className="matrix-bg" ref={matrixRef}></div>
      
      <div className="container">
        {showSuccess && <div className="success-anim"></div>}
        
        <h1 
          ref={headingRef}
          className={`glitch-text ${isTextGlitching ? 'intense' : ''}`} 
          data-text="ENTER YOUR EMAIL"
          onMouseOver={() => setIsTextGlitching(true)}
          onMouseOut={() => setIsTextGlitching(false)}
        >
          ENTER YOUR EMAIL
        </h1>
        
        {!submitted ? (
          <form onSubmit={handleSubmit} className="glitch-container">
            <div className="glitch-item">
              <input
                type="email"
                value={email}
                onChange={handleEmailChange}
                placeholder="your.email@example.com"
                required
                className="glitch-green glitch-item"
                maxLength={maxChars}
              />
              <div style={{ textAlign: 'right', fontSize: '0.8rem', color: charCount > 80 ? 'var(--neon-pink)' : 'var(--neon-green)' }}>
                {charCount}/{maxChars}
              </div>
            </div>
            
            <div className="glitch-item">
              <button 
                type="submit" 
                className="glitch-pink glitch-item"
                onMouseOver={() => setIsTextGlitching(true)}
                onMouseOut={() => setIsTextGlitching(false)}
              >
                SUBMIT
              </button>
            </div>
          </form>
        ) : (
          <div className="glitch-container">
            <p className="glitch-text" data-text="WELCOME TO THE">WELCOME TO THE</p>
            <p className="glitch-text" data-text="FUTURE">FUTURE</p>
            <p 
              ref={emailRef}
              className="glitch-text email-glitch" 
              data-text={maskedEmail}
            >
              {maskedEmail}
            </p>
            <p className="glitch-text" data-text="CONNECTED">CONNECTED</p>
          </div>
        )}
      </div>
    </>
  );
}
