data "cwdashboard_metric" "this" {
  metric_name = "CPUUtilization"
  namespace   = "AWS/EC2"
  dimensions_map = {
    InstanceId = "i-0123456789abcdef0"
  }
  statistic = "Average"
  period    = 60
}

data "cwdashboard_graph_widget" "this" {
  width  = 24
  height = 6

  title = "EC2 CPU Utilization"

  left = [
    data.cwdashboard_metric.this.json,
  ]
}

data "cwdashboard" "this" {
  start           = "-PT7D"
  period_override = "auto"
  widgets = [
    data.cwdashboard_graph_widget.this.json,
  ]
}

# use dashboard JSON with the Terraform AWS Provider
resource "aws_cloudwatch_dashboard" "this" {
  dashboard_name = "test-dashboard"
  dashboard_body = data.cwdashboard.this.json
}
