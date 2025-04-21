// ========================================
// Internet Gateway & Route Table for NAT
// ========================================
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.cv_c9_vpc.id

  tags = {
    Name = "${var.project}-igw"
  }
}

