  // Sample chart drawing function
        function drawBarChart() {
            const canvas = document.getElementById('barChart');
            const ctx = canvas.getContext('2d');
            const data = document.getElementById('barData').value.split(',').map(Number);
            const threshold = parseInt(document.getElementById('barThreshold').value);
            const title = document.getElementById('chartTitle').value;
            
            // Update chart title
            document.getElementById('dynamicChartTitle').textContent = title;
            
            // Set canvas dimensions
            canvas.width = Math.max(600, data.length * 60);
            canvas.height = 300;
            
            // Clear canvas
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            
            // Chart settings
            const barWidth = 40;
            const barSpacing = 60;
            const chartHeight = 250;
            const chartTop = 30;
            const maxValue = Math.max(...data, threshold) * 1.1;
            
            // Draw bars
            data.forEach((value, index) => {
                const barHeight = (value / maxValue) * chartHeight;
                const x = index * barSpacing + 30;
                const y = chartTop + chartHeight - barHeight;
                
                // Bar color based on threshold
                ctx.fillStyle = value >= threshold ? '#28a745' : '#dc3545';
                ctx.fillRect(x, y, barWidth, barHeight);
                
                // Value labels
                ctx.fillStyle = '#333';
                ctx.font = '12px Arial';
                ctx.textAlign = 'center';
                ctx.fillText(value, x + barWidth/2, y - 5);
                
                // Index labels
                ctx.fillText(index + 1, x + barWidth/2, chartTop + chartHeight + 20);
            });
            
            // Draw threshold line
            const thresholdY = chartTop + chartHeight - (threshold / maxValue) * chartHeight;
            ctx.strokeStyle = '#ff6b35';
            ctx.lineWidth = 2;
            ctx.setLineDash([5, 5]);
            ctx.beginPath();
            ctx.moveTo(20, thresholdY);
            ctx.lineTo(canvas.width - 20, thresholdY);
            ctx.stroke();
            
            // Threshold label
            ctx.fillStyle = '#ff6b35';
            ctx.font = 'bold 12px Arial';
            ctx.textAlign = 'left';
            ctx.fillText(`Batas: ${threshold}`, 25, thresholdY - 5);
        }
        
        function loadSampleData1() {
            document.getElementById('barData').value = '180,220,195,240,160,210,185,200';
            document.getElementById('chartTitle').value = 'Data Produksi Minggu 1';
            drawBarChart();
        }
        
        function loadSampleData2() {
            document.getElementById('barData').value = '150,175,190,165,140,185,170,195';
            document.getElementById('chartTitle').value = 'Data Produksi Minggu 2';
            drawBarChart();
        }
        
        function loadSampleData3() {
            document.getElementById('barData').value = '200,180,220,195,175,205,190,185';
            document.getElementById('chartTitle').value = 'Data Produksi Minggu 3';
            drawBarChart();
        }
        
        function clearChart() {
            document.getElementById('barData').value = '';
            document.getElementById('chartTitle').value = 'Grafik Produksi';
            document.getElementById('barThreshold').value = '150';
            const canvas = document.getElementById('barChart');
            const ctx = canvas.getContext('2d');
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            document.getElementById('dynamicChartTitle').textContent = 'Grafik Produksi';
        }
        
        function exportChart() {
            const canvas = document.getElementById('barChart');
            const link = document.createElement('a');
            link.download = 'chart.png';
            link.href = canvas.toDataURL();
            link.click();
        }
        
        // Initialize chart on page load
        window.onload = function() {
            drawBarChart();
        };