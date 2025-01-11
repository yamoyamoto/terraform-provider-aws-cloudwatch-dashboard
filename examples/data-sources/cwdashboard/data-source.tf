data "cwdashboard_text_widget" "this" {
  markdown   = "Hello, World!"
  background = "#000000"
  width      = 24
  height     = 2
}

data "cwdashboard_metric" "this" {
  metric_name = "CPUUtilization"
  namespace   = "AWS/EC2"
  dimensions_map = {
    InstanceId = "i-0123456789abcdef0"
  }
  color     = "#0000ff"
  statistic = "Average"
  period    = 60
}

data "cwdashboard_graph_widget" "this" {
  width  = 24
  height = 6
  title  = "EC2 CPU Utilization"

  left = [
    data.cwdashboard_metric.this.json,
  ]
}

data "cwdashboard" "this" {
  start           = "-PT7D"
  period_override = "auto"
  widgets = [
    data.cwdashboard_text_widget.this.json,
    data.cwdashboard_graph_widget.this.json,
  ]
}
