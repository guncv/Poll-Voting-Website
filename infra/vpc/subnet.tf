// ========================================
// Public Subnet for NAT Gateway
// ========================================
resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.cv_c9_vpc.id
  cidr_block              = "10.0.3.0/24"
  availability_zone       = var.availability_zones[0]  # AZ[0] - First availability zone
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project}-public-subnet"
  }

  depends_on = [aws_vpc.cv_c9_vpc]
}


resource "aws_subnet" "public_2" {
  vpc_id                  = aws_vpc.cv_c9_vpc.id
  cidr_block              = "10.0.4.0/24"
  availability_zone       = var.availability_zones[1] # Ensure this is in a different AZ
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project}-public-subnet-2"
  }

  depends_on = [aws_vpc.cv_c9_vpc]
}

resource "aws_subnet" "public_3" {
  vpc_id                  = aws_vpc.cv_c9_vpc.id
  cidr_block              = "10.0.5.0/24"
  availability_zone       = var.availability_zones[2]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project}-public-subnet-3"
  }

  depends_on = [aws_vpc.cv_c9_vpc]
}


// ========================================
// Private Subnets
// ========================================

resource "aws_subnet" "private_1" {
  vpc_id                  = aws_vpc.cv_c9_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone = var.availability_zones[0]
  map_public_ip_on_launch = false

  tags = {
    Name = "${var.project}-private-subnet-1"
  }

  depends_on = [aws_vpc.cv_c9_vpc]
}

resource "aws_subnet" "private_2" {
  vpc_id                  = aws_vpc.cv_c9_vpc.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = var.availability_zones[1]
  map_public_ip_on_launch = false

  tags = {
    Name = "${var.project}-private-subnet-2"
  }

  depends_on = [aws_vpc.cv_c9_vpc]
}
